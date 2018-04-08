package ledge

import (
	"github.com/msaldanha/realChain/blockstore"
	"github.com/msaldanha/realChain/block"
	"errors"
	"github.com/msaldanha/realChain/keypair"
	"github.com/msaldanha/realChain/address"
)

type Ledge struct {
	bs *blockstore.BlockStore
	accounts map[string]*Account
}

func New() (*Ledge) {
	return &Ledge{accounts: make(map[string]*Account, 0)}
}

func (ld *Ledge) Use(bs *blockstore.BlockStore) {
	ld.bs = bs
}

func (ld *Ledge) Initialize(initialBalance float64) (*block.Block, error) {
	if !ld.bs.IsEmpty() {
		return nil, errors.New("Ledge already initialized")
	}

	genesisBlock := ld.bs.CreateOpenBlock()
	acc, err := ld.CreateAccount()
	if err != nil {
		return nil, err
	}
	genesisBlock.Account = []byte(acc)
	genesisBlock.Representative = genesisBlock.Account
	genesisBlock.Balance = initialBalance

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
			return nil, errors.New("Ledge already initialized")
		}
		return nil, err
	}
	return blk, nil
}

func (ld *Ledge) Send(from, to string, amount float64) (string, error) {
	fromTipBlk, err := ld.bs.Retrieve(from)

	if err != nil {
		return "", err
	}

	if from == to {
		return "", errors.New("from and to accounts must be different")
	}

	if fromTipBlk == nil {
		return "", errors.New("from account not found")
	}

	if fromTipBlk.Balance < amount {
		return "", errors.New("not enough funds")
	}

	if amount == 0 {
		return "", errors.New("amount cannot be zero")
	}

	hash, err := ld.createSendTransaction(fromTipBlk, []byte(to), amount)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (ld *Ledge) Receive(send *block.Block) (string, error) {
	acc := ld.GetAccount(send.Link)

	if acc == nil {
		return "", errors.New("account not managed by this ledge")
	}

	hash, err := ld.createReceiveTransaction(send)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (ld *Ledge) createSendTransaction(fromTip *block.Block, to []byte, amount float64) ([]byte, error) {
	send := ld.bs.CreateSendBlock()
	send.Account = fromTip.Account
	send.Link = to
	send.Previous = fromTip.Hash
	send.Balance = fromTip.Balance - amount
	if err := ld.signAndPow(send); err != nil {
		return nil, err
	}
	send, err := ld.bs.Store(send)
	if err != nil {
		return nil, err
	}
	return send.Hash, nil
}

func (ld *Ledge) createReceiveTransaction(send *block.Block) ([]byte, error) {
	prev, err := ld.bs.Retrieve(string(send.Previous))
	if err != nil {
		return nil, err
	}

	amount := prev.Balance - send.Balance
	if amount <= 0 {
		return nil, errors.New("invalid account balance")
	}

	addr := address.New()
	if valid, err := addr.IsValid(string(send.Link)); !valid {
		return nil, errors.New("send transaction addr is not valid: " + err.Error())
	}

	acc := ld.GetAccount(send.Link)
	if acc == nil {
		return nil, errors.New("account not managed by this ledge")
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
	} else {
		receive = ld.bs.CreateOpenBlock()
		receive.Balance = amount
		receive.Representative = send.Link
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

func (ld *Ledge) signAndPow(blk *block.Block) (error) {
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

func (ld *Ledge) CreateAccount() (string, error) {
	keys, err := keypair.New()
	if err != nil {
		return "", err
	}

	acc := &Account{Keys: keys}
	addr := address.New()
	ad, err := addr.GenerateForKey(acc.Keys.PublicKey)
	if err != nil {
		return "", err
	}

	acc.Address = string(ad)
	ld.AddAccount(acc)
	return acc.Address, nil
}

func (ld *Ledge) AddAccount(acc *Account) {
	ld.accounts[acc.Address] = acc
}

func (ld *Ledge) GetDefaultAccount() []byte {
	return []byte("account")
}

func (ld *Ledge) GetAccount(acc []byte) *Account {
	return ld.accounts[string(acc)]
}

func (ld *Ledge) setPow(blk *block.Block) error {
	nonce, hash, err := ld.bs.CalculatePow(blk)
	if err != nil {
		return err
	}
	blk.PowNonce = nonce
	blk.Hash = hash
	return nil
}

func (ld *Ledge) sign(blk *block.Block) error {
	hash, err := blk.GetHash()
	if err != nil {
		return err
	}
	blk.Signature = hash
	return nil
}

//TODO: add transaction verification as in the paper