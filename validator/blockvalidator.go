package validator

import (
	. "github.com/msaldanha/realChain/block"
	"errors"
	"github.com/msaldanha/realChain/keyvaluestore"
	"strconv"
	"strings"
	"crypto/sha256"
	"encoding/hex"
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


func (bs *BaseBlockValidator) CreateHash(block *Block) (string, error) {
	parts := [...]string{
		strconv.FormatInt(block.Timestamp, 10),
		strconv.Itoa(int(block.Type)),
		block.Account,
		block.Representative,
		block.Previous,
		block.Link,
		block.Work,
		strconv.FormatFloat(block.Balance, 'f', -1, 64),
	}
	all := strings.Join(parts[:], "|")
	allBytes := []byte(all)
	sh := sha256.Sum256(allBytes)
	return hex.EncodeToString(sh[:]), nil
}

func (bs *BaseBlockValidator) HasValidSignature(block *Block) (bool, error) {
	hash, _ := bs.CreateHash(block)
	if block.Signature == hash {
		return true, nil
	}
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
	if block.Previous == "" {
		return false, errors.New("Previous block can not be empty")
	}
	if block.Signature == "" {
		return false, errors.New("Block signature can not be empty")
	}
	if block.Work == "" {
		return false, errors.New("Block PoW can not be empty")
	}
	return true, nil
}

func (v *OpenBlockValidator) IsFilled(block *Block) (bool, error) {
	if !block.Type.IsValid() || block.Type != OPEN {
		return false, errors.New("Invalid block type")
	}
	if block.Account == "" {
		return false, errors.New("Block account can not be empty")
	}
	if block.Representative == "" {
		return false, errors.New("Block representative can not be empty")
	}
	if block.Signature == "" {
		return false, errors.New("Block signature can not be empty")
	}
	if block.Work == "" {
		return false, errors.New("Block PoW can not be empty")
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
	if block.Link == "" {
		return false, errors.New("Block destination can not be empty")
	}
	return v.BaseBlockValidator.IsFilled(block)
}

func (v *SendBlockValidator) IsValid(block *Block) (bool, error) {
	if ok, err := v.IsFilled(block); !ok {
		return ok, err
	}
	if ok, err := v.HasValidDestination(block); !ok {
		return ok, err
	}
	return v.BaseBlockValidator.IsValid(block)
}

func (v *SendBlockValidator) HasValidDestination(block *Block) (bool, error) {
	_, found, err := v.store.Get(block.Link)
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
	if block.Link == "" {
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
	dest, found, err := v.store.Get(block.Link)
	if err != nil {
		return false, err
	}
	if !found {
		return false, errors.New("Source not found")
	}
	source := dest.(*Block)
	if source.Type != SEND {
		return false, errors.New("Source of invalid type")
	}
	return true, nil
}

func (v *ChangeBlockValidator) IsFilled(block *Block) (bool, error) {
	if !block.Type.IsValid() || block.Type != CHANGE {
		return false, errors.New("Invalid block type")
	}
	if block.Representative == "" {
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