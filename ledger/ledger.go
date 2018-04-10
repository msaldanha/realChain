package ledger

import (
	"github.com/msaldanha/realChain/blockstore"
	"github.com/msaldanha/realChain/block"
	"github.com/msaldanha/realChain/keypair"
	"github.com/msaldanha/realChain/address"
	"github.com/msaldanha/realChain/Error"
	"crypto/ecdsa"
	"crypto/rand"
	"math/big"
	"crypto/elliptic"
	"bytes"
)

const (
	ErrInvalidAccountBalance         = Error.Error("invalid account balance")
	ErrLedgerAlreadyInitialized      = Error.Error("ledger already initialized")
	ErrFromToMustBeDifferent         = Error.Error("from and to accounts must be different")
	ErrFromAccountNotFound           = Error.Error("from account not found")
	ErrNotEnoughFunds                = Error.Error("not enough funds")
	ErrAmountCantBeZero              = Error.Error("amount cannot be zero")
	ErrAccountNotManagedByThisLedger = Error.Error("account not managed by this ledger")
	ErrInvalidSendTransactionAddress = Error.Error("invalid send transaction address")
	ErrInvalidTransactionSignature   = Error.Error("invalid transaction signature")
	ErrInvalidTransactionHash        = Error.Error("invalid transaction hash")
	ErrTransactionAlreadyInLedger    = Error.Error("transaction already in ledger")
	ErrTransactionNotFound           = Error.Error("previous transaction not found")
	ErrPreviousTransactionNotFound   = Error.Error("previous transaction not found")
	ErrHeadTransactionNotFound       = Error.Error("head transaction not found")
	ErrPreviousTransactionIsNotHead  = Error.Error("previous transaction is not the chain head")
	ErrOpenTransactionNotFound       = Error.Error("open transaction not found")
	ErrAccountDoesNotMatchPubKey     = Error.Error("account does not match public key")
)

type Ledger struct {
	bs       *blockstore.BlockStore
	accounts map[string]*Account
}

func New() (*Ledger) {
	return &Ledger{accounts: make(map[string]*Account, 0)}
}

func (ld *Ledger) Use(bs *blockstore.BlockStore) {
	ld.bs = bs
}

func (ld *Ledger) Initialize(initialBalance float64) (*block.Block, error) {
	if !ld.bs.IsEmpty() {
		return nil, ErrLedgerAlreadyInitialized
	}

	genesisBlock := ld.bs.CreateOpenBlock()
	acc, err := ld.CreateAccount()
	if err != nil {
		return nil, err
	}
	genesisBlock.Account = []byte(acc.Address)
	genesisBlock.Representative = genesisBlock.Account
	genesisBlock.Balance = initialBalance
	genesisBlock.PubKey = acc.Keys.PublicKey

	err = ld.setPow(genesisBlock)
	if err != nil {
		return nil, err
	}

	err = ld.sign(genesisBlock)
	if err != nil {
		return nil, err
	}

	blk, err := ld.bs.Store(genesisBlock)
	if err != nil {
		if err.Error() == "Previous block can not be empty" {
			return nil, ErrLedgerAlreadyInitialized
		}
		return nil, err
	}
	return blk, nil
}

func (ld *Ledger) Send(from, to string, amount float64) (string, error) {
	fromTipBlk, err := ld.bs.Retrieve(from)

	if err != nil {
		return "", err
	}

	addr := address.New()
	if valid, err := addr.IsValid(to); !valid {
		return "", err
	}

	if from == to {
		return "", ErrFromToMustBeDifferent
	}

	if fromTipBlk == nil {
		return "", ErrFromAccountNotFound
	}

	if fromTipBlk.Balance < amount {
		return "", ErrNotEnoughFunds
	}

	if amount == 0 {
		return "", ErrAmountCantBeZero
	}

	hash, err := ld.createSendTransaction(fromTipBlk, []byte(to), amount)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (ld *Ledger) Receive(send *block.Block) (string, error) {
	acc := ld.GetAccount(send.Link)

	if acc == nil {
		return "", ErrAccountNotManagedByThisLedger
	}

	hash, err := ld.createReceiveTransaction(send)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (ld *Ledger) createSendTransaction(fromTip *block.Block, to []byte, amount float64) ([]byte, error) {
	send := ld.bs.CreateSendBlock()
	send.Account = fromTip.Account
	send.Link = to
	send.Previous = fromTip.Hash
	send.Balance = fromTip.Balance - amount
	send.PubKey = fromTip.PubKey
	if err := ld.signAndPow(send); err != nil {
		return nil, err
	}
	send, err := ld.bs.Store(send)
	if err != nil {
		return nil, err
	}
	return send.Hash, nil
}

func (ld *Ledger) createReceiveTransaction(send *block.Block) ([]byte, error) {
	prev, err := ld.bs.Retrieve(string(send.Previous))
	if err != nil {
		return nil, err
	}

	if ok, err := ld.verifyTransaction(send, false); !ok {
		return nil, err
	}

	amount := prev.Balance - send.Balance
	if amount <= 0 {
		return nil, ErrInvalidAccountBalance
	}

	addr := address.New()
	if valid, _ := addr.IsValid(string(send.Link)); !valid {
		return nil, ErrInvalidSendTransactionAddress
	}

	acc := ld.GetAccount(send.Link)
	if acc == nil {
		return nil, ErrAccountNotManagedByThisLedger
	}

	receiveTip, err := ld.bs.Retrieve(string(send.Link))
	if err != nil {
		return nil, err
	}

	var receive *block.Block
	if receiveTip != nil {
		receive = ld.bs.CreateReceiveBlock()
		receive.Previous = receiveTip.Hash
		receive.Balance = receiveTip.Balance + amount
		receive.Representative = receiveTip.Representative
		receive.PubKey = receiveTip.PubKey
	} else {
		receive = ld.bs.CreateOpenBlock()
		receive.Balance = amount
		receive.Representative = send.Link
		receive.PubKey = acc.Keys.PublicKey
	}

	receive.Account = send.Link
	receive.Link = send.Hash

	if err := ld.signAndPow(receive); err != nil {
		return nil, err
	}

	receive, err = ld.bs.Store(receive)
	if err != nil {
		return nil, err
	}

	return receive.Hash, nil
}

func (ld *Ledger) signAndPow(blk *block.Block) (error) {
	err := ld.setPow(blk)
	if err != nil {
		return err
	}

	err = ld.sign(blk)
	if err != nil {
		return err
	}

	return nil
}

func (ld *Ledger) CreateAccount() (*Account, error) {
	keys, err := keypair.New()
	if err != nil {
		return nil, err
	}

	acc := &Account{Keys: keys}
	addr := address.New()
	ad, err := addr.GenerateForKey(acc.Keys.PublicKey)
	if err != nil {
		return nil, err
	}

	acc.Address = string(ad)
	ld.AddAccount(acc)
	return acc, nil
}

func (ld *Ledger) AddAccount(acc *Account) {
	ld.accounts[acc.Address] = acc
}

func (ld *Ledger) GetAccount(acc []byte) *Account {
	return ld.accounts[string(acc)]
}

func (ld *Ledger) setPow(blk *block.Block) error {
	nonce, hash, err := ld.bs.CalculatePow(blk)
	if err != nil {
		return err
	}
	blk.PowNonce = nonce
	blk.Hash = hash
	return nil
}

func (ld *Ledger) sign(blk *block.Block) error {
	signature, err := ld.getSignature(blk)
	if err != nil {
		return err
	}
	blk.Signature = signature
	return nil
}

func (ld *Ledger) getSignature(blk *block.Block) ([]byte, error) {
	acc := ld.GetAccount(blk.Account)
	if acc == nil {
		return nil, ErrAccountNotManagedByThisLedger
	}
	privateKey := ld.getPrivateKey(acc)
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, blk.Hash)
	if err != nil {
		return nil, err
	}
	return append(r.Bytes(), s.Bytes()...), nil
}

func (ld *Ledger) verifySignature(blk *block.Block) bool {
	r := big.Int{}
	s := big.Int{}
	sigLen := len(blk.Signature)
	r.SetBytes(blk.Signature[:(sigLen / 2)])
	s.SetBytes(blk.Signature[(sigLen / 2):])

	x := big.Int{}
	y := big.Int{}
	keyLen := len(blk.PubKey)
	x.SetBytes(blk.PubKey[:(keyLen / 2)])
	y.SetBytes(blk.PubKey[(keyLen / 2):])

	curve := elliptic.P256()
	rawPubKey := ecdsa.PublicKey{curve, &x, &y}

	return ecdsa.Verify(&rawPubKey, blk.Hash, &r, &s)
}

func (ld *Ledger) verifyPow(blk *block.Block) bool {
	ok, _ := ld.bs.VerifyPow(blk)
	return ok
}

func (ld *Ledger) getPrivateKey(acc *Account) *ecdsa.PrivateKey {
	D := new(big.Int)
	D.SetBytes(acc.Keys.PrivateKey)

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

func (ld *Ledger) verifyTransaction(blk *block.Block, isNew bool) (bool, error) {
	if ok, err := ld.verifyAddress(blk); !ok {
		return false, err
	}
	if !ld.verifyPow(blk) {
		return false, ErrInvalidTransactionHash
	}
	if !ld.verifySignature(blk) {
		return false, ErrInvalidTransactionSignature
	}

	b, err := ld.bs.Retrieve(string(blk.Hash))
	if err != nil {
		return false, err
	}
	if b != nil && isNew {
		return false, ErrTransactionAlreadyInLedger
	} else if b == nil && !isNew {
		return false, ErrTransactionNotFound
	}

	if blk.Type != block.OPEN {
		previous, err := ld.bs.Retrieve(string(blk.Previous))
		if err != nil {
			return false, err
		}
		if previous == nil {
			return false, ErrPreviousTransactionNotFound
		}
		if isNew {
			head, err := ld.bs.Retrieve(string(previous.Account))
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

	open, _ := ld.getOpenTransaction(blk)
	if open == nil {
		return false, ErrOpenTransactionNotFound
	}
	return true, nil
}

func (ld *Ledger) verifyAddress(blk *block.Block) (bool, error) {
	addr := address.New()
	ad, err := addr.GenerateForKey(blk.PubKey)
	if err != nil {
		return false, err
	}
	if ad != string(blk.Account) {
		return false, ErrAccountDoesNotMatchPubKey
	}
	return true, nil
}

func (ld *Ledger) getOpenTransaction(blk *block.Block) (*block.Block, error) {
	var ret = blk
	var err error
	for ret != nil && ret.Type != block.OPEN {
		ret, err = ld.getPreviousTransaction(ret)
		if err != nil {
			return nil, err
		}
	}
	if ret != nil && ret.Type == block.OPEN {
		return ret, nil
	}
	return nil, nil
}

func (ld *Ledger) getPreviousTransaction(blk *block.Block) (*block.Block, error) {
	previous, err := ld.bs.Retrieve(string(blk.Previous))
	if err != nil {
		return nil, err
	}
	if previous == nil {
		return nil, ErrPreviousTransactionNotFound
	}
	return previous, nil
}
