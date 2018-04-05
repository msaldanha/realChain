package blockstore

import (
	. "github.com/msaldanha/realChain/block"
	. "github.com/msaldanha/realChain/keyvaluestore"
	"errors"
	"github.com/msaldanha/realChain/validator"
)

type BlockStore struct {
	store                 Storer
	blockValidatorCreator validator.BlockValidatorCreator
}

func New(store Storer, validatorCreator validator.BlockValidatorCreator) (*BlockStore) {
	a := &BlockStore{store: store, blockValidatorCreator: validatorCreator}
	return a
}

func (bs *BlockStore) isValid(block *Block) (bool, error) {
	if !block.Type.IsValid(){
		return false, errors.New("Invalid block type")
	}
	val := bs.blockValidatorCreator.CreateValidatorForBlock(block.Type, bs.store)
	return val.IsValid(block)
}

func (bs *BlockStore) Store(block *Block) (*Block, error) {
	if ok, err := bs.isValid(block); !ok {
		return nil, err
	}
	bs.store.Put(string(block.Hash), block)
	return block, nil
}

func (bs *BlockStore) Retrieve(hash string) (*Block, error) {
	value, found, err := bs.store.Get(hash)
	if err != nil {
		return nil, err
	}

	if found {
		if blk, ok := value.(*Block); ok {
			return blk, nil
		}
		return nil, errors.New("Not a block")
	}

	return nil, nil
}
