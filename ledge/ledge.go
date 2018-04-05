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