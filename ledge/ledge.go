package ledge

import (
	"github.com/msaldanha/realChain/blockstore"
	"github.com/msaldanha/realChain/block"
	"errors"
)

type Ledge struct {
	bs *blockstore.BlockStore
}

func New() (*Ledge) {
	return &Ledge{}
}

func (ld *Ledge) Use(bs *blockstore.BlockStore) {
	ld.bs = bs
}

func (ld *Ledge) Initialize(initialBalance float64) (*block.Block, error) {
	genesisBlock := ld.bs.CreateOpenBlock()
	genesisBlock.Link = []byte("Genesis")
	genesisBlock.Account = ld.CreateAccount()
	genesisBlock.Representative = genesisBlock.Account
	genesisBlock.Balance = initialBalance

	err := ld.setPow(genesisBlock)
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

	if fromTipBlk == nil {
		return "", errors.New("from account not found")
	}

	if fromTipBlk.Balance < amount {
		return "", errors.New("not enough funds")
	}

	hash, err := ld.createSendTransaction(fromTipBlk, []byte(to), amount)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (ld *Ledge) createSendTransaction(fromTip *block.Block, to []byte, amount float64) ([]byte, error) {
	send := ld.bs.CreateSendBlock()
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

func (ld *Ledge) CreateAccount() []byte {
	return []byte("account")
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
