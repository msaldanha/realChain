package ledger

import (
	"github.com/msaldanha/realChain/transactionstore"
	"github.com/msaldanha/realChain/transaction"
	"github.com/msaldanha/realChain/keypair"
	"github.com/msaldanha/realChain/address"
	"github.com/msaldanha/realChain/Error"
	"crypto/ecdsa"
	"crypto/rand"
	"math/big"
	"crypto/elliptic"
	"bytes"
	"github.com/msaldanha/realChain/keyvaluestore"
	log "github.com/sirupsen/logrus"
	"math"
)

const (
	ErrInvalidAddressBalance               = Error.Error("invalid address balance")
	ErrLedgerAlreadyInitialized            = Error.Error("ledger already initialized")
	ErrFromToMustBeDifferent               = Error.Error("from and to addresses must be different")
	ErrFromAddressNotFound                 = Error.Error("from address not found")
	ErrNotEnoughFunds                      = Error.Error("not enough funds")
	ErrAmountCantBeZero                    = Error.Error("amount cannot be zero")
	ErrAddressNotManagedByThisLedger       = Error.Error("address not managed by this ledger")
	ErrInvalidSendTransactionAddress       = Error.Error("invalid send transaction address")
	ErrInvalidTransactionSignature         = Error.Error("invalid transaction signature")
	ErrInvalidTransactionHash              = Error.Error("invalid transaction hash")
	ErrTransactionAlreadyInLedger          = Error.Error("transaction already in ledger")
	ErrTransactionNotFound                 = Error.Error("previous transaction not found")
	ErrPreviousTransactionNotFound         = Error.Error("previous transaction not found")
	ErrHeadTransactionNotFound             = Error.Error("head transaction not found")
	ErrPreviousTransactionIsNotHead        = Error.Error("previous transaction is not the chain head")
	ErrSendTransactionIsNotHead            = Error.Error("send transaction is not the chain head")
	ErrSendTransactionIsNotPending         = Error.Error("send transaction is not pending")
	ErrOpenTransactionNotFound             = Error.Error("open transaction not found")
	ErrAddressDoesNotMatchPubKey           = Error.Error("address does not match public key")
	ErrSendTransactionNotFound             = Error.Error("send transaction not found")
	ErrSentAmountDiffersFromReceivedAmount = Error.Error("sent amount differs from received amount")
)

type Ledger struct {
	ts        *transactionstore.TransactionStore
	addresses keyvaluestore.Storer
}

func New() (*Ledger) {
	return &Ledger{}
}

func (ld *Ledger) Use(txStore *transactionstore.TransactionStore, addressStore keyvaluestore.Storer) {
	ld.ts = txStore
	ld.addresses = addressStore
}

func (ld *Ledger) Initialize(initialBalance float64) (*transaction.Transaction, error) {
	if !ld.ts.IsEmpty() {
		return nil, ErrLedgerAlreadyInitialized
	}

	genesisTx := ld.ts.CreateOpenTransaction()
	addr, err := ld.CreateAddress()
	if err != nil {
		return nil, err
	}
	genesisTx.Address = []byte(addr.Address)
	genesisTx.Representative = genesisTx.Address
	genesisTx.Balance = initialBalance
	genesisTx.PubKey = addr.Keys.PublicKey

	err = ld.setPow(genesisTx)
	if err != nil {
		return nil, err
	}

	err = ld.sign(genesisTx)
	if err != nil {
		return nil, err
	}

	_, err = ld.saveTransaction(genesisTx)
	if err != nil {
		if err.Error() == "Previous transaction can not be empty" {
			return nil, ErrLedgerAlreadyInitialized
		}
		return nil, err
	}
	return genesisTx, nil
}

func (ld *Ledger) CreateSendTransaction(from, to string, amount float64) (*transaction.Transaction, error) {
	fromTipTx, err := ld.ts.Retrieve(from)

	if err != nil {
		return nil, err
	}

	addr := address.New()
	if valid, err := addr.IsValid(to); !valid {
		return nil, err
	}

	if from == to {
		return nil, ErrFromToMustBeDifferent
	}

	if fromTipTx == nil {
		return nil, ErrFromAddressNotFound
	}

	if fromTipTx.Balance < amount {
		return nil, ErrNotEnoughFunds
	}

	if amount == 0 {
		return nil, ErrAmountCantBeZero
	}

	tx, err := ld.createSendTransaction(fromTipTx, []byte(to), amount)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (ld *Ledger) Receive(send *transaction.Transaction) (string, error) {
	addr, err := ld.GetAddress(send.Link)
	if err != nil {
		logError("Receive", send, err)
		return "", err
	}
	if addr == nil {
		logError("Receive", send, ErrAddressNotManagedByThisLedger)
		return "", ErrAddressNotManagedByThisLedger
	}

	hash, err := ld.createReceiveTransaction(send)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (ld *Ledger) HandleTransactionBytes(txBytes []byte) (error) {
	tx := transaction.NewTransactionFromBytes(txBytes)
	return ld.HandleTransaction(tx)
}

func (ld *Ledger) HandleTransaction(tx *transaction.Transaction) (err error) {
	err = nil
	if ok, err := ld.verifyTransaction(tx, true); !ok {
		logError("HandleTransaction", tx, err)
		return err
	}

	_, err = ld.saveTransaction(tx)
	if err != nil {
		logError("HandleTransaction", tx, err)
		return
	}

	if tx.Type == transaction.SEND {
		addr, err := ld.GetAddress(tx.Link)
		if err != nil {
			return err
		}
		if addr != nil {
			_, err = ld.createReceiveTransaction(tx)
		}
	}

	return
}

func (ld *Ledger) GetLastTransaction(address string) (*transaction.Transaction, error) {
	fromTipTx, err := ld.ts.Retrieve(address)
	if err != nil {
		return nil, err
	}
	return fromTipTx, nil
}

func (ld *Ledger) GetAddressStatement(address string) ([]*transaction.Transaction, error) {
	txChain, err := ld.ts.GetTransactionChain(address, false)
	if err != nil {
		return nil, err
	}
	return txChain, nil
}

func (ld *Ledger) createSendTransaction(fromTip *transaction.Transaction, to []byte, amount float64) (*transaction.Transaction, error) {
	send := ld.ts.CreateSendTransaction()
	send.Address = fromTip.Address
	send.Link = to
	send.Previous = fromTip.Hash
	send.Balance = fromTip.Balance - amount
	send.PubKey = fromTip.PubKey
	if err := ld.signAndPow(send); err != nil {
		return nil, err
	}
	return send, nil
}

func (ld *Ledger) createReceiveTransaction(send *transaction.Transaction) ([]byte, error) {
	prev, err := ld.ts.Retrieve(string(send.Previous))
	if err != nil {
		return nil, err
	}

	if ok, err := ld.verifyTransaction(send, false); !ok {
		return nil, err
	}

	amount := prev.Balance - send.Balance
	if amount <= 0 {
		return nil, ErrInvalidAddressBalance
	}

	addr := address.New()
	if valid, _ := addr.IsValid(string(send.Link)); !valid {
		return nil, ErrInvalidSendTransactionAddress
	}

	addr1, err := ld.GetAddress(send.Link)
	if err != nil {
		return nil, err
	}
	if addr1 == nil {
		return nil, ErrAddressNotManagedByThisLedger
	}

	receiveTip, err := ld.ts.Retrieve(string(send.Link))
	if err != nil {
		return nil, err
	}

	var receive *transaction.Transaction
	if receiveTip != nil {
		receive = ld.ts.CreateReceiveTransaction()
		receive.Previous = receiveTip.Hash
		receive.Balance = receiveTip.Balance + amount
		receive.Representative = receiveTip.Representative
		receive.PubKey = receiveTip.PubKey
	} else {
		receive = ld.ts.CreateOpenTransaction()
		receive.Balance = amount
		receive.Representative = send.Link
		receive.PubKey = addr1.Keys.PublicKey
	}

	receive.Address = send.Link
	receive.Link = send.Hash

	if err := ld.signAndPow(receive); err != nil {
		return nil, err
	}

	err = ld.HandleTransaction(receive)
	if err != nil {
		return nil, err
	}

	return receive.Hash, nil
}

func (ld *Ledger) saveTransaction(tx *transaction.Transaction) ([]byte, error) {
	tx, err := ld.ts.Store(tx)
	if err != nil {
		return nil, err
	}
	return tx.Hash, nil
}

func (ld *Ledger) signAndPow(tx *transaction.Transaction) (error) {
	err := ld.setPow(tx)
	if err != nil {
		return err
	}

	err = ld.sign(tx)
	if err != nil {
		return err
	}

	return nil
}

func (ld *Ledger) CreateAddress() (*Address, error) {
	keys, err := keypair.New()
	if err != nil {
		return nil, err
	}

	addr1 := &Address{Keys: keys}
	addr := address.New()
	ad, err := addr.GenerateForKey(addr1.Keys.PublicKey)
	if err != nil {
		return nil, err
	}

	addr1.Address = string(ad)
	err = ld.AddAddress(addr1)
	if err != nil {
		return nil, err
	}
	return addr1, nil
}

func (ld *Ledger) AddAddress(addr *Address) error {
	return ld.addresses.Put(addr.Address, addr.ToBytes())
}

func (ld *Ledger) GetAddress(address []byte) (*Address, error) {
	addr, ok, err := ld.addresses.Get(string(address))
	if !ok {
		return nil, err
	}
	return NewAddressFromBytes(addr), nil
}

func (ld *Ledger) GetAddresses() ([]*Address, error) {
	addrs, err := ld.addresses.GetAll()
	if err != nil {
		return nil, err
	}
	addresses := make([]*Address, 0)
	for _, v := range addrs {
		addresses = append(addresses, NewAddressFromBytes(v))
	}
	return addresses, nil
}

func (ld *Ledger) setPow(tx *transaction.Transaction) error {
	nonce, hash, err := ld.ts.CalculatePow(tx)
	if err != nil {
		return err
	}
	tx.PowNonce = nonce
	tx.Hash = hash
	return nil
}

func (ld *Ledger) sign(tx *transaction.Transaction) error {
	signature, err := ld.getSignature(tx)
	if err != nil {
		return err
	}
	tx.Signature = signature
	return nil
}

func (ld *Ledger) getSignature(tx *transaction.Transaction) ([]byte, error) {
	addr, err := ld.GetAddress(tx.Address)
	if err != nil {
		return nil, err
	}
	if addr == nil {
		return nil, ErrAddressNotManagedByThisLedger
	}
	privateKey := ld.getPrivateKey(addr)
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, tx.Hash)
	if err != nil {
		return nil, err
	}
	return append(r.Bytes(), s.Bytes()...), nil
}

func (ld *Ledger) verifySignature(tx *transaction.Transaction) bool {
	r := big.Int{}
	s := big.Int{}
	sigLen := len(tx.Signature)
	r.SetBytes(tx.Signature[:(sigLen / 2)])
	s.SetBytes(tx.Signature[(sigLen / 2):])

	x := big.Int{}
	y := big.Int{}
	keyLen := len(tx.PubKey)
	x.SetBytes(tx.PubKey[:(keyLen / 2)])
	y.SetBytes(tx.PubKey[(keyLen / 2):])

	curve := elliptic.P256()
	rawPubKey := ecdsa.PublicKey{curve, &x, &y}

	return ecdsa.Verify(&rawPubKey, tx.Hash, &r, &s)
}

func (ld *Ledger) verifyPow(tx *transaction.Transaction) bool {
	ok, _ := ld.ts.VerifyPow(tx)
	return ok
}

func (ld *Ledger) getPrivateKey(addr *Address) *ecdsa.PrivateKey {
	D := new(big.Int)
	D.SetBytes(addr.Keys.PrivateKey)

	curve := elliptic.P256()
	privateKey := ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: curve,
		},
		D: D,
	}

	privateKey.PublicKey.X, privateKey.PublicKey.Y = curve.ScalarBaseMult(D.Bytes())
	return &privateKey
}

func (ld *Ledger) verifyTransaction(tx *transaction.Transaction, isNew bool) (bool, error) {
	if ok, err := ld.verifyAddress(tx); !ok {
		return false, err
	}
	if !ld.verifyPow(tx) {
		return false, ErrInvalidTransactionHash
	}
	if !ld.verifySignature(tx) {
		return false, ErrInvalidTransactionSignature
	}

	b, err := ld.ts.Retrieve(string(tx.Hash))
	if err != nil {
		return false, err
	}
	if b != nil && isNew {
		return false, ErrTransactionAlreadyInLedger
	} else if b == nil && !isNew {
		return false, ErrTransactionNotFound
	}

	if tx.Type != transaction.OPEN {
		previous, err := ld.ts.Retrieve(string(tx.Previous))
		if err != nil {
			return false, err
		}
		if previous == nil {
			return false, ErrPreviousTransactionNotFound
		}
		if isNew {
			head, err := ld.ts.Retrieve(string(previous.Address))
			if err != nil {
				return false, err
			}
			if head == nil {
				return false, ErrHeadTransactionNotFound
			}
			if bytes.Compare(head.Hash, previous.Hash) != 0 {
				return false, ErrPreviousTransactionIsNotHead
			}
		}
	}

	if tx.Type == transaction.OPEN || tx.Type == transaction.RECEIVE {
		send, err := ld.ts.Retrieve(string(tx.Link))
		if err != nil {
			return false, err
		}
		if send == nil {
			return false, ErrSendTransactionNotFound
		}
		if isNew {
			pending, err := ld.isPendingTransaction(send)
			if err != nil {
				return false, err
			}
			if !pending {
				return false, ErrSendTransactionIsNotPending
			}
			head, err := ld.ts.Retrieve(string(send.Address))
			if err != nil {
				return false, err
			}
			if head == nil {
				return false, ErrHeadTransactionNotFound
			}
			if bytes.Compare(head.Hash, send.Hash) != 0 {
				return false, ErrSendTransactionIsNotHead
			}
			sentAmount, err := ld.findAbsoluteBalanceDiffWithPrevious(send)
			if err != nil {
				return false, err
			}
			receivedAmount, err := ld.findAbsoluteBalanceDiffWithPrevious(tx)
			if err != nil {
				return false, err
			}
			if sentAmount != receivedAmount {
				return false, ErrSentAmountDiffersFromReceivedAmount
			}
		}
	}

	open, _ := ld.getOpenTransaction(tx)
	if open == nil {
		return false, ErrOpenTransactionNotFound
	}
	return true, nil
}

func (ld *Ledger) isPendingTransaction(tx *transaction.Transaction) (bool, error) {
	if tx.Type != transaction.SEND {
		return false, nil
	}
	target, err := ld.GetLastTransaction(string(tx.Link))
	if err != nil {
		return false, err
	}
	if target == nil {
		return true, nil
	}
	chain, err := ld.ts.GetTransactionChain(string(target.Hash), false)
	if err != nil {
		return false, err
	}
	for _, v := range chain {
		if bytes.Equal(tx.Hash, v.Link) {
			return false, nil
		}
	}
	return true, nil
}

func (ld *Ledger) findAbsoluteBalanceDiffWithPrevious(tx *transaction.Transaction) (float64, error) {
	var amount float64 = 0
	previous, err := ld.findPrevious(tx)
	if err != nil {
		return 0, err
	}
	if previous == nil {
		return 0, err
	}
	if tx.Type == transaction.OPEN {
		amount = tx.Balance
	} else {
		amount = tx.Balance - previous.Balance
	}
	return math.Abs(amount), nil
}

func (ld *Ledger) findPrevious(tx *transaction.Transaction) (*transaction.Transaction, error) {
	if tx.Type != transaction.OPEN {
		return ld.ts.Retrieve(string(tx.Previous))
	}
	return ld.ts.Retrieve(string(tx.Link))
}

func (ld *Ledger) verifyAddress(tx *transaction.Transaction) (bool, error) {
	addr := address.New()
	ad, err := addr.GenerateForKey(tx.PubKey)
	if err != nil {
		return false, err
	}
	if ad != string(tx.Address) {
		return false, ErrAddressDoesNotMatchPubKey
	}
	return true, nil
}

func (ld *Ledger) getOpenTransaction(tx *transaction.Transaction) (*transaction.Transaction, error) {
	var ret = tx
	var err error
	for ret != nil && ret.Type != transaction.OPEN {
		ret, err = ld.getPreviousTransaction(ret)
		if err != nil {
			return nil, err
		}
	}
	if ret != nil && ret.Type == transaction.OPEN {
		return ret, nil
	}
	return nil, nil
}

func (ld *Ledger) getPreviousTransaction(tx *transaction.Transaction) (*transaction.Transaction, error) {
	previous, err := ld.ts.Retrieve(string(tx.Previous))
	if err != nil {
		return nil, err
	}
	if previous == nil {
		return nil, ErrPreviousTransactionNotFound
	}
	return previous, nil
}

func logError(context string, tx *transaction.Transaction, err error) {
	log.Printf("Ledger.%s failed: %s (tx: %s)", context, err, string(tx.Hash))
}