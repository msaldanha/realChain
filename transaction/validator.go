package transaction

import (
	"github.com/msaldanha/realChain/keyvaluestore"
	"github.com/msaldanha/realChain/Error"
)

const (
	ErrInvalidTransactionType               = Error.Error("invalid transaction type")
	ErrInvalidTransactionTimestamp          = Error.Error("invalid transaction timestamp")
	ErrTransactionAddressCantBeEmpty        = Error.Error("transaction address can not be empty")
	ErrPreviousTransactionCantBeEmpty       = Error.Error("previous transaction can not be empty")
	ErrTransactionSignatureCantBeEmpty      = Error.Error("transaction signature can not be empty")
	ErrTransactionPowNonceCantBeZero        = Error.Error("transaction PoW nonce can not be zero")
	ErrTransactionHashCantBeEmpty           = Error.Error("transaction hash can not be empty")
	ErrTransactionLinkCantBeEmpty           = Error.Error("transaction link can not be empty")
	ErrTransactionRepresentativeCantBeEmpty = Error.Error("transaction representative can not be empty")
	ErrDestinationNotFound                  = Error.Error("destination not found")
	ErrSourceNotFound                       = Error.Error("source not found")
	ErrInvalidSourceType                    = Error.Error("invalid source type")
	ErrPubKeyCantBeEmpty                    = Error.Error("transaction public key can not be empty")
)

type Validator interface {
	IsFilled(tx *Transaction) (bool, error)
	IsValid(tx *Transaction) (bool, error)
}

type ValidatorCreator interface {
	CreateValidatorForTransaction(txType Type, store keyvaluestore.Storer) (Validator)
}

type validatorCreator struct {
}

type BaseValidator struct {
	store keyvaluestore.Storer
}

type OpenValidator struct {
	BaseValidator
}

type SendValidator struct {
	BaseValidator
}

type ReceiveValidator struct {
	BaseValidator
}

type ChangeValidator struct {
	BaseValidator
}

func NewValidatorCreator() (ValidatorCreator) {
	return &validatorCreator{}
}

func (*validatorCreator) CreateValidatorForTransaction(txType Type, store keyvaluestore.Storer) (Validator) {
	switch txType {
	case OPEN:
		return &OpenValidator{BaseValidator{store}}
	case SEND:
		return &SendValidator{BaseValidator{store}}
	case RECEIVE:
		return &ReceiveValidator{BaseValidator{store}}
	case CHANGE:
		return &ChangeValidator{BaseValidator{store}}
	}
	return nil
}

func (v *BaseValidator) IsValid(tx *Transaction) (bool, error) {
	return v.IsFilled(tx)
}

func (v *BaseValidator) IsFilled(tx *Transaction) (bool, error) {
	if !tx.Type.IsValid() {
		return false, ErrInvalidTransactionType
	}
	if tx.Timestamp <= 0 {
		return false, ErrInvalidTransactionTimestamp
	}
	if len(tx.Address) == 0 {
		return false, ErrTransactionAddressCantBeEmpty
	}
	if len(tx.Previous) == 0 && !v.store.IsEmpty() && tx.Type != OPEN {
		return false, ErrPreviousTransactionCantBeEmpty
	}
	if len(tx.Signature) == 0 {
		return false, ErrTransactionSignatureCantBeEmpty
	}
	if tx.PowNonce == 0 {
		return false, ErrTransactionPowNonceCantBeZero
	}
	if len(tx.Hash) == 0 {
		return false, ErrTransactionHashCantBeEmpty
	}
	if len(tx.PubKey) == 0 {
		return false, ErrPubKeyCantBeEmpty
	}
	return true, nil
}

func (v *OpenValidator) IsFilled(tx *Transaction) (bool, error) {
	if !tx.Type.IsValid() || tx.Type != OPEN {
		return false, ErrInvalidTransactionType
	}
	if ok, err := v.BaseValidator.IsFilled(tx); !ok {
		return ok, err
	}
	if len(tx.Link) == 0 && !v.store.IsEmpty() {
		return false, ErrTransactionLinkCantBeEmpty
	}
	if len(tx.Representative) == 0 {
		return false, ErrTransactionRepresentativeCantBeEmpty
	}
	if len(tx.Signature) == 0 {
		return false, ErrTransactionSignatureCantBeEmpty
	}
	if tx.PowNonce == 0 {
		return false, ErrTransactionPowNonceCantBeZero
	}
	return true, nil
}

func (v *OpenValidator) IsValid(tx *Transaction) (bool, error) {
	if ok, err := v.IsFilled(tx); !ok {
		return ok, err
	}
	return v.BaseValidator.IsValid(tx)
}

func (v *SendValidator) IsFilled(tx *Transaction) (bool, error) {
	if !tx.Type.IsValid() || tx.Type != SEND {
		return false, ErrInvalidTransactionType
	}
	if ok, err := v.BaseValidator.IsFilled(tx); !ok {
		return ok, err
	}
	if len(tx.Link) == 0 {
		return false, ErrTransactionLinkCantBeEmpty
	}
	return v.BaseValidator.IsFilled(tx)
}

func (v *SendValidator) IsValid(tx *Transaction) (bool, error) {
	if ok, err := v.IsFilled(tx); !ok {
		return ok, err
	}
	return v.BaseValidator.IsValid(tx)
}

func (v *SendValidator) HasValidDestination(tx *Transaction) (bool, error) {
	_, found, err := v.store.Get(string(tx.Link))
	if err != nil {
		return false, err
	}
	if !found {
		return false, ErrDestinationNotFound
	}
	return found, err
}

func (v *ReceiveValidator) IsFilled(tx *Transaction) (bool, error) {
	if !tx.Type.IsValid() || tx.Type != RECEIVE {
		return false, ErrInvalidTransactionType
	}
	if ok, err := v.BaseValidator.IsFilled(tx); !ok {
		return ok, err
	}
	if len(tx.Link) == 0 {
		return false, ErrTransactionLinkCantBeEmpty
	}
	return v.BaseValidator.IsFilled(tx)
}

func (v *ReceiveValidator) IsValid(tx *Transaction) (bool, error) {
	if ok, err := v.IsFilled(tx); !ok {
		return ok, err
	}
	if ok, err := v.HasValidSource(tx); !ok {
		return ok, err
	}
	return v.BaseValidator.IsValid(tx)
}

func (v *ReceiveValidator) HasValidSource(tx *Transaction) (bool, error) {
	dest, found, err := v.store.Get(string(tx.Link))
	if err != nil {
		return false, err
	}
	if !found {
		return false, ErrSourceNotFound
	}
	source := NewTransactionFromBytes(dest)
	if source.Type != SEND {
		return false, ErrInvalidSourceType
	}
	return true, nil
}

func (v *ChangeValidator) IsFilled(tx *Transaction) (bool, error) {
	if !tx.Type.IsValid() || tx.Type != CHANGE {
		return false, ErrInvalidTransactionType
	}
	if ok, err := v.BaseValidator.IsFilled(tx); !ok {
		return ok, err
	}
	if len(tx.Representative) == 0 {
		return false, ErrTransactionRepresentativeCantBeEmpty
	}
	return v.BaseValidator.IsFilled(tx)
}

func (v *ChangeValidator) IsValid(tx *Transaction) (bool, error) {
	if ok, err := v.IsFilled(tx); !ok {
		return ok, err
	}
	return v.BaseValidator.IsValid(tx)
}
