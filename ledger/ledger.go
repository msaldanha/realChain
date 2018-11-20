package ledger

import (
	"github.com/msaldanha/realChain/transaction"
	"github.com/msaldanha/realChain/Error"
	"github.com/msaldanha/realChain/address"
)

const (
	ErrInvalidOperation                    = Error.Error("invalid operation")
	ErrInvalidAddressBalance               = Error.Error("invalid address balance")
	ErrLedgerAlreadyInitialized            = Error.Error("ledger already initialized")
	ErrFromToMustBeDifferent               = Error.Error("from and to addresses must be different")
	ErrFromAddressNotFound                 = Error.Error("from address not found")
	ErrNotEnoughFunds                      = Error.Error("not enough funds")
	ErrAmountCantBeZero                    = Error.Error("amount cannot be zero")
	ErrAddressNotManagedByThisLedger       = Error.Error("address not managed by this ledger")
	ErrInvalidSendTransactionAddress       = Error.Error("invalid send transaction address")
	ErrInvalidTransactionSignature         = Error.Error("invalid transaction signature")
	ErrInvalidTransactionHash              = Error.Error("invalid transaction hash")
	ErrTransactionAlreadyInLedger          = Error.Error("transaction already in ledger")
	ErrTransactionNotFound                 = Error.Error("previous transaction not found")
	ErrPreviousTransactionNotFound         = Error.Error("previous transaction not found")
	ErrHeadTransactionNotFound             = Error.Error("head transaction not found")
	ErrPreviousTransactionIsNotHead        = Error.Error("previous transaction is not the chain head")
	ErrSendTransactionIsNotHead            = Error.Error("send transaction is not the chain head")
	ErrSendTransactionIsNotPending         = Error.Error("send transaction is not pending")
	ErrOpenTransactionNotFound             = Error.Error("open transaction not found")
	ErrAddressDoesNotMatchPubKey           = Error.Error("address does not match public key")
	ErrSendTransactionNotFound             = Error.Error("send transaction not found")
	ErrSentAmountDiffersFromReceivedAmount = Error.Error("sent amount differs from received amount")
)

//go:generate mockgen -destination=../tests/mock_ledger.go -package=tests github.com/msaldanha/realChain/ledger Ledger

type Ledger interface {
	Initialize(initialBalance float64) (*transaction.Transaction, *address.Address, error)
	GetLastTransaction(address string) (*transaction.Transaction, error)
	GetTransaction(hash string) (*transaction.Transaction, error)
	GetAddressStatement(address string) ([]*transaction.Transaction, error)
	HandleTransactionBytes(txBytes []byte) (*transaction.Transaction, error)
	HandleTransaction(tx *transaction.Transaction) (ret *transaction.Transaction, err error)
	AddAddress(addr *address.Address) error
	Receive(send *transaction.Transaction) (string, error)
	VerifyTransaction(tx *transaction.Transaction, isNew bool) error
}
