package block

import (
	"github.com/msaldanha/realChain/keyvaluestore"
	"bytes"
	"github.com/msaldanha/realChain/Error"
)

const (
	ErrBlockSignatureDoesNotMatch     = Error.Error("block signature does not match")
	ErrInvalidBlockType               = Error.Error("invalid block type")
	ErrInvalidBlockTimestamp          = Error.Error("invalid block timestamp")
	ErrBlockAccountCantBeEmpty        = Error.Error("block account can not be empty")
	ErrPreviousBlockCantBeEmpty       = Error.Error("previous block can not be empty")
	ErrBlockSignatureCantBeEmpty      = Error.Error("block signature can not be empty")
	ErrBlockPowNonceCantBeZero        = Error.Error("block PoW nonce can not be zero")
	ErrBlockHashCantBeEmpty           = Error.Error("block hash can not be empty")
	ErrBlockLinkCantBeEmpty           = Error.Error("block link can not be empty")
	ErrBlockRepresentativeCantBeEmpty = Error.Error("block representative can not be empty")
	ErrDestinationNotFound            = Error.Error("destination not found")
	ErrSourceNotFound                 = Error.Error("source not found")
	ErrInvalidSourceType              = Error.Error("invalid source type")
)

type Validator interface {
	IsFilled(block *Block) (bool, error)
	IsValid(block *Block) (bool, error)
}

type ValidatorCreator interface {
	CreateValidatorForBlock(blockType BlockType, store keyvaluestore.Storer) (Validator)
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

func NewBlockValidatorCreator() (ValidatorCreator) {
	return &blockValidatorCreator{}
}

func (*blockValidatorCreator) CreateValidatorForBlock(blockType BlockType, store keyvaluestore.Storer) (Validator) {
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
	return false, ErrBlockSignatureDoesNotMatch
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
		return false, ErrInvalidBlockType
	}
	if block.Timestamp <= 0 {
		return false, ErrInvalidBlockTimestamp
	}
	if len(block.Account) == 0 {
		return false, ErrBlockAccountCantBeEmpty
	}
	if len(block.Previous) == 0 && !v.store.IsEmpty() && block.Type != OPEN {
		return false, ErrPreviousBlockCantBeEmpty
	}
	if len(block.Signature) == 0 {
		return false, ErrBlockSignatureCantBeEmpty
	}
	if block.PowNonce == 0 {
		return false, ErrBlockPowNonceCantBeZero
	}
	if len(block.Hash) == 0 {
		return false, ErrBlockHashCantBeEmpty
	}
	return true, nil
}

func (v *OpenBlockValidator) IsFilled(block *Block) (bool, error) {
	if !block.Type.IsValid() || block.Type != OPEN {
		return false, ErrInvalidBlockType
	}
	if ok, err := v.BaseBlockValidator.IsFilled(block); !ok {
		return ok, err
	}
	if len(block.Link) == 0 && !v.store.IsEmpty() {
		return false, ErrBlockLinkCantBeEmpty
	}
	if len(block.Representative) == 0 {
		return false, ErrBlockRepresentativeCantBeEmpty
	}
	if len(block.Signature) == 0 {
		return false, ErrBlockSignatureCantBeEmpty
	}
	if block.PowNonce == 0 {
		return false, ErrBlockPowNonceCantBeZero
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
		return false, ErrInvalidBlockType
	}
	if ok, err := v.BaseBlockValidator.IsFilled(block); !ok {
		return ok, err
	}
	if len(block.Link) == 0 {
		return false, ErrBlockLinkCantBeEmpty
	}
	return v.BaseBlockValidator.IsFilled(block)
}

func (v *SendBlockValidator) HasValidDestination(block *Block) (bool, error) {
	_, found, err := v.store.Get(string(block.Link))
	if err != nil {
		return false, err
	}
	if !found {
		return false, ErrDestinationNotFound
	}
	return found, err
}

func (v *ReceiveBlockValidator) IsFilled(block *Block) (bool, error) {
	if !block.Type.IsValid() || block.Type != RECEIVE {
		return false, ErrInvalidBlockType
	}
	if ok, err := v.BaseBlockValidator.IsFilled(block); !ok {
		return ok, err
	}
	if len(block.Link) == 0 {
		return false, ErrBlockLinkCantBeEmpty
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
		return false, ErrSourceNotFound
	}
	source := dest.(*Block)
	if source.Type != SEND {
		return false, ErrInvalidSourceType
	}
	return true, nil
}

func (v *ChangeBlockValidator) IsFilled(block *Block) (bool, error) {
	if !block.Type.IsValid() || block.Type != CHANGE {
		return false, ErrInvalidBlockType
	}
	if ok, err := v.BaseBlockValidator.IsFilled(block); !ok {
		return ok, err
	}
	if len(block.Representative) == 0 {
		return false, ErrBlockRepresentativeCantBeEmpty
	}
	return v.BaseBlockValidator.IsFilled(block)
}

func (v *ChangeBlockValidator) IsValid(block *Block) (bool, error) {
	if ok, err := v.IsFilled(block); !ok {
		return ok, err
	}
	return v.BaseBlockValidator.IsValid(block)
}
