package validator

import (
	. "github.com/msaldanha/realChain/block"
	"errors"
	"github.com/msaldanha/realChain/keyvaluestore"
	"bytes"
)

type BlockValidator interface {
	IsFilled(block *Block) (bool, error)
	IsValid(block *Block) (bool, error)
}

type BlockValidatorCreator interface {
	CreateValidatorForBlock(blockType BlockType, store keyvaluestore.Storer) (BlockValidator)
}

type blockValidatorCreator struct {
}

type BaseBlockValidator struct {
	store keyvaluestore.Storer
}

type OpenBlockValidator struct {
	BaseBlockValidator
}

type SendBlockValidator struct {
	BaseBlockValidator
}

type ReceiveBlockValidator struct {
	BaseBlockValidator
}

type ChangeBlockValidator struct {
	BaseBlockValidator
}

func NewBlockValidatorCreator() (BlockValidatorCreator) {
	return &blockValidatorCreator{}
}

func (*blockValidatorCreator) CreateValidatorForBlock(blockType BlockType, store keyvaluestore.Storer) (BlockValidator) {
	switch blockType {
	case OPEN:
		return &OpenBlockValidator{BaseBlockValidator{store}}
	case SEND:
		return &SendBlockValidator{BaseBlockValidator{store}}
	case RECEIVE:
		return &ReceiveBlockValidator{BaseBlockValidator{store}}
	case CHANGE:
		return &ChangeBlockValidator{BaseBlockValidator{store}}
	}
	return nil
}

func (v *BaseBlockValidator) HasValidSignature(block *Block) (bool, error) {
	hash, _ := block.GetHash()
	if bytes.Compare(block.Signature, hash) == 0 {
		return true, nil
	}
	//text := fmt.Sprintf("%s : %s", string(block.Signature), string(hash))
	//log.Println(text)
	return false, errors.New("Block signature does not match")
}

func (v *BaseBlockValidator) IsValid(block *Block) (bool, error) {
	if ok, err := v.IsFilled(block); !ok {
		return ok, err
	}
	if ok, err := v.HasValidSignature(block); !ok {
		return ok, err
	}
	return true, nil
}

func (v *BaseBlockValidator) IsFilled(block *Block) (bool, error) {
	if !block.Type.IsValid() {
		return false, errors.New("Invalid block type")
	}
	if block.Timestamp <= 0 {
		return false, errors.New("Invalid block timestamp")
	}
	if len(block.Previous) == 0 && !v.store.IsEmpty() {
		return false, errors.New("Previous block can not be empty")
	}
	if len(block.Signature) == 0 {
		return false, errors.New("Block signature can not be empty")
	}
	if block.PowNonce == 0 {
		return false, errors.New("Block PoW nonce can not be zero")
	}
	if len(block.Hash) == 0 {
		return false, errors.New("Block hash can not be empty")
	}
	return true, nil
}

func (v *OpenBlockValidator) IsFilled(block *Block) (bool, error) {
	if !block.Type.IsValid() || block.Type != OPEN {
		return false, errors.New("Invalid block type")
	}
	if len(block.Account) == 0 {
		return false, errors.New("Block account can not be empty")
	}
	if len(block.Representative) == 0 {
		return false, errors.New("Block representative can not be empty")
	}
	if len(block.Signature) == 0 {
		return false, errors.New("Block signature can not be empty")
	}
	if block.PowNonce == 0 {
		return false, errors.New("Block PoW nonce can not be zero")
	}
	return true, nil
}

func (v *OpenBlockValidator) IsValid(block *Block) (bool, error) {
	if ok, err := v.IsFilled(block); !ok {
		return ok, err
	}
	return v.BaseBlockValidator.IsValid(block)
}

func (v *SendBlockValidator) IsFilled(block *Block) (bool, error) {
	if !block.Type.IsValid() || block.Type != SEND {
		return false, errors.New("Invalid block type")
	}
	if len(block.Link) == 0 {
		return false, errors.New("Block destination can not be empty")
	}
	return v.BaseBlockValidator.IsFilled(block)
}

func (v *SendBlockValidator) HasValidDestination(block *Block) (bool, error) {
	_, found, err := v.store.Get(string(block.Link))
	if err != nil {
		return false, err
	}
	if !found {
		return false, errors.New("Destination not found")
	}
	return found, err
}

func (v *ReceiveBlockValidator) IsFilled(block *Block) (bool, error) {
	if !block.Type.IsValid() || block.Type != RECEIVE {
		return false, errors.New("Invalid block type")
	}
	if len(block.Link) == 0 {
		return false, errors.New("Block source can not be empty")
	}
	return v.BaseBlockValidator.IsFilled(block)
}

func (v *ReceiveBlockValidator) IsValid(block *Block) (bool, error) {
	if ok, err := v.IsFilled(block); !ok {
		return ok, err
	}
	if ok, err := v.HasValidSource(block); !ok {
		return ok, err
	}
	return v.BaseBlockValidator.IsValid(block)
}

func (v *ReceiveBlockValidator) HasValidSource(block *Block) (bool, error) {
	dest, found, err := v.store.Get(string(block.Link))
	if err != nil {
		return false, err
	}
	if !found {
		return false, errors.New("Source not found")
	}
	source := dest
	if source.Type != SEND {
		return false, errors.New("Source of invalid type")
	}
	return true, nil
}

func (v *ChangeBlockValidator) IsFilled(block *Block) (bool, error) {
	if !block.Type.IsValid() || block.Type != CHANGE {
		return false, errors.New("Invalid block type")
	}
	if len(block.Representative) == 0 {
		return false, errors.New("Block representative can not be empty")
	}
	return v.BaseBlockValidator.IsFilled(block)
}

func (v *ChangeBlockValidator) IsValid(block *Block) (bool, error) {
	if ok, err := v.IsFilled(block); !ok {
		return ok, err
	}
	return v.BaseBlockValidator.IsValid(block)
}