package ledger

import (
	"github.com/msaldanha/realChain/transactionstore"
	"github.com/msaldanha/realChain/transaction"
	"github.com/msaldanha/realChain/address"
	"crypto/ecdsa"
	"bytes"
	"github.com/msaldanha/realChain/keyvaluestore"
	log "github.com/sirupsen/logrus"
	"math"
)

type LocalLedger struct {
	ts        *transactionstore.TransactionStore
	addresses keyvaluestore.Storer
}

func NewLocalLedger(txStore *transactionstore.TransactionStore, addressStore keyvaluestore.Storer) (*LocalLedger) {
	return &LocalLedger{ts:txStore, addresses:addressStore}
}

func (ld *LocalLedger) Initialize(initialBalance float64) (*transaction.Transaction, *address.Address, error) {
	if !ld.ts.IsEmpty() {
		return nil, nil, ErrLedgerAlreadyInitialized
	}

	genesisTx := transaction.NewOpenTransaction()
	addr, err := ld.CreateAddress()
	if err != nil {
		return nil, nil, err
	}

	ld.AddAddress(addr)

	genesisTx.Address = []byte(addr.Address)
	genesisTx.Representative = genesisTx.Address
	genesisTx.Balance = initialBalance
	genesisTx.PubKey = addr.Keys.PublicKey

	err = ld.setPow(genesisTx)
	if err != nil {
		return nil, nil, err
	}

	err = ld.sign(genesisTx)
	if err != nil {
		return nil, nil, err
	}

	_, err = ld.saveTransaction(genesisTx)
	if err != nil {
		if err.Error() == "Previous transaction can not be empty" {
			return nil, nil, ErrLedgerAlreadyInitialized
		}
		return nil, nil, err
	}
	return genesisTx, addr, nil
}

func (ld *LocalLedger) CreateSendTransaction(from, to string, amount float64) (*transaction.Transaction, error) {
	fromTipTx, err := ld.ts.Retrieve(from)

	if err != nil {
		return nil, err
	}

	addr := address.New()
	addr.Address = to
	if valid, err := addr.IsValid(); !valid {
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

func (ld *LocalLedger) Receive(send *transaction.Transaction) (string, error) {
	if err := ld.VerifyTransaction(send, false); err != nil {
		logError("HandleTransaction", send, err)
		return "", err
	}
	addr, err := ld.GetAddress(send.Link)
	if err != nil {
		logError("Receive", send, err)
		return "", err
	}
	if addr == nil {
		logError("Receive", send, ErrAddressNotManagedByThisLedger)
		return "", ErrAddressNotManagedByThisLedger
	}

	tx, err := ld.createReceiveTransaction(send)
	if err != nil {
		return "", err
	}

	return string(tx.Hash), nil
}

func (ld *LocalLedger) HandleTransactionBytes(txBytes []byte) (*transaction.Transaction, error) {
	tx := transaction.NewTransactionFromBytes(txBytes)
	return ld.HandleTransaction(tx)
}

func (ld *LocalLedger) HandleTransaction(tx *transaction.Transaction) (ret *transaction.Transaction, err error) {
	err = nil
	ret = nil
	if err = ld.VerifyTransaction(tx, true); err != nil {
		logError("HandleTransaction", tx, err)
		return
	}

	_, err = ld.saveTransaction(tx)
	if err != nil {
		logError("HandleTransaction", tx, err)
		return
	}

	if tx.Type == transaction.SEND {
		addr, err := ld.GetAddress(tx.Link)
		if err != nil {
			return nil, err
		}
		if addr != nil {
			return ld.createReceiveTransaction(tx)
		}
	}

	return
}

func (ld *LocalLedger) GetLastTransaction(address string) (*transaction.Transaction, error) {
	fromTipTx, err := ld.ts.Retrieve(address)
	if err != nil {
		return nil, err
	}
	return fromTipTx, nil
}

func (ld *LocalLedger) GetTransaction(hash string) (*transaction.Transaction, error) {
	tx, err := ld.ts.Retrieve(hash)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (ld *LocalLedger) GetAddressStatement(address string) ([]*transaction.Transaction, error) {
	txChain, err := ld.ts.GetTransactionChain(address, false)
	if err != nil {
		return nil, err
	}
	return txChain, nil
}

func (ld *LocalLedger) createSendTransaction(fromTip *transaction.Transaction, to []byte, amount float64) (*transaction.Transaction, error) {
	send := transaction.NewSendTransaction()
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

func (ld *LocalLedger) createReceiveTransaction(send *transaction.Transaction) (*transaction.Transaction, error) {
	prev, err := ld.ts.Retrieve(string(send.Previous))
	if err != nil {
		return nil, err
	}

	if err := ld.VerifyTransaction(send, false); err != nil {
		return nil, err
	}

	amount := prev.Balance - send.Balance
	if amount <= 0 {
		return nil, ErrInvalidAddressBalance
	}

	addr := address.New()
	addr.Address = string(send.Link)
	if valid, _ := addr.IsValid(); !valid {
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
		receive = transaction.NewReceiveTransaction()
		receive.Previous = receiveTip.Hash
		receive.Balance = receiveTip.Balance + amount
		receive.Representative = receiveTip.Representative
		receive.PubKey = receiveTip.PubKey
	} else {
		receive = transaction.NewOpenTransaction()
		receive.Balance = amount
		receive.Representative = send.Link
		receive.PubKey = addr1.Keys.PublicKey
	}

	receive.Address = send.Link
	receive.Link = send.Hash

	if err := ld.signAndPow(receive); err != nil {
		return nil, err
	}

	_, err = ld.HandleTransaction(receive)
	if err != nil {
		return nil, err
	}

	return receive, nil
}

func (ld *LocalLedger) saveTransaction(tx *transaction.Transaction) ([]byte, error) {
	tx, err := ld.ts.Store(tx)
	if err != nil {
		return nil, err
	}
	return tx.Hash, nil
}

func (ld *LocalLedger) signAndPow(tx *transaction.Transaction) (error) {
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

func (ld *LocalLedger) CreateAddress() (*address.Address, error) {
	return address.NewAddressWithKeys()
}

func (ld *LocalLedger) AddAddress(addr *address.Address) error {
	return ld.addresses.Put(addr.Address, addr.ToBytes())
}

func (ld *LocalLedger) GetAddress(addressBytes []byte) (*address.Address, error) {
	addr, ok, err := ld.addresses.Get(string(addressBytes))
	if !ok {
		return nil, err
	}
	return address.NewAddressFromBytes(addr), nil
}

func (ld *LocalLedger) GetAddresses() ([]*address.Address, error) {
	addrs, err := ld.addresses.GetAll()
	if err != nil {
		return nil, err
	}
	addresses := make([]*address.Address, 0)
	for _, v := range addrs {
		addresses = append(addresses, address.NewAddressFromBytes(v))
	}
	return addresses, nil
}

func (ld *LocalLedger) setPow(tx *transaction.Transaction) error {
	nonce, hash, err := tx.CalculatePow()
	if err != nil {
		return err
	}
	tx.PowNonce = nonce
	tx.Hash = hash
	return nil
}

func (ld *LocalLedger) sign(tx *transaction.Transaction) error {
	addr, err := ld.GetAddress(tx.Address)
	if err != nil {
		return err
	}
	if addr == nil {
		return ErrAddressNotManagedByThisLedger
	}
	privateKey := ld.getPrivateKey(addr)
	err = tx.Sign(privateKey)
	return err
}

func (ld *LocalLedger) verifyPow(tx *transaction.Transaction) bool {
	ok, _ := tx.VerifyPow()
	return ok
}

func (ld *LocalLedger) getPrivateKey(addr *address.Address) *ecdsa.PrivateKey {
	return addr.Keys.ToEcdsaPrivateKey()
}

func (ld *LocalLedger) VerifyTransaction(tx *transaction.Transaction, isNew bool) (error) {
	if ok, err := ld.verifyAddress(tx); !ok {
		return err
	}
	if ok, err := ld.verifyLinkAddress(tx); !ok {
		return err
	}
	if !ld.verifyPow(tx) {
		return ErrInvalidTransactionHash
	}
	if !tx.VerifySignature() {
		return ErrInvalidTransactionSignature
	}

	localTx, err := ld.ts.Retrieve(string(tx.Hash))
	if err != nil {
		return err
	}
	if localTx != nil && isNew {
		return ErrTransactionAlreadyInLedger
	} else if localTx == nil && !isNew {
		return ErrTransactionNotFound
	}

	if tx.Type != transaction.OPEN {
		previous, err := ld.ts.Retrieve(string(tx.Previous))
		if err != nil {
			return err
		}
		if previous == nil {
			return ErrPreviousTransactionNotFound
		}
		if isNew {
			head, err := ld.ts.Retrieve(string(previous.Address))
			if err != nil {
				return err
			}
			if head == nil {
				return ErrHeadTransactionNotFound
			}
			if bytes.Compare(head.Hash, previous.Hash) != 0 {
				return ErrPreviousTransactionIsNotHead
			}
		}
	}

	if tx.Type == transaction.SEND {
		amount, err := ld.findBalanceDiffWithPrevious(tx)
		if err != nil {
			return err
		}
		if amount < 0.0 || tx.Balance < 0.0 {
			return ErrNotEnoughFunds
		}
	}

	if tx.Type == transaction.OPEN || tx.Type == transaction.RECEIVE {
		send, err := ld.ts.Retrieve(string(tx.Link))
		if err != nil {
			return err
		}
		if send == nil {
			return ErrSendTransactionNotFound
		}
		if isNew {
			pending, err := ld.isPendingTransaction(send)
			if err != nil {
				return err
			}
			if !pending {
				return ErrSendTransactionIsNotPending
			}
			head, err := ld.ts.Retrieve(string(send.Address))
			if err != nil {
				return err
			}
			if head == nil {
				return ErrHeadTransactionNotFound
			}
			if bytes.Compare(head.Hash, send.Hash) != 0 {
				return ErrSendTransactionIsNotHead
			}
			sentAmount, err := ld.findAbsoluteBalanceDiffWithPrevious(send)
			if err != nil {
				return err
			}
			receivedAmount, err := ld.findAbsoluteBalanceDiffWithPrevious(tx)
			if err != nil {
				return err
			}
			if sentAmount != receivedAmount {
				return ErrSentAmountDiffersFromReceivedAmount
			}
		}
	}

	open, _ := ld.getOpenTransaction(tx)
	if open == nil {
		return ErrOpenTransactionNotFound
	}
	return nil
}

func (ld *LocalLedger) isPendingTransaction(tx *transaction.Transaction) (bool, error) {
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

func (ld *LocalLedger) findAbsoluteBalanceDiffWithPrevious(tx *transaction.Transaction) (float64, error) {
	amount, err := ld.findBalanceDiffWithPrevious(tx)
	return math.Abs(amount), err
}

func (ld *LocalLedger) findBalanceDiffWithPrevious(tx *transaction.Transaction) (float64, error) {
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
		amount = previous.Balance - tx.Balance
	}
	return amount, nil
}

func (ld *LocalLedger) findPrevious(tx *transaction.Transaction) (*transaction.Transaction, error) {
	if tx.Type != transaction.OPEN {
		return ld.ts.Retrieve(string(tx.Previous))
	}
	return ld.ts.Retrieve(string(tx.Link))
}

func (ld *LocalLedger) verifyAddress(tx *transaction.Transaction) (bool, error) {
	if !address.MatchesPubKey(tx.Address, tx.PubKey) {
		return false, ErrAddressDoesNotMatchPubKey
	}
	return true, nil
}

func (ld *LocalLedger) verifyLinkAddress(tx *transaction.Transaction) (bool, error) {
	if tx.Type == transaction.SEND {
		return address.IsValid(string(tx.Link))
	}
	return true, nil
}

func (ld *LocalLedger) getOpenTransaction(tx *transaction.Transaction) (*transaction.Transaction, error) {
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

func (ld *LocalLedger) getPreviousTransaction(tx *transaction.Transaction) (*transaction.Transaction, error) {
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