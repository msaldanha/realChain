package ledger

import (
	"github.com/msaldanha/realChain/Error"
	"github.com/msaldanha/realChain/transaction"
)

const (
	ErrLedgerAlreadyInitialized                 = Error.Error("ledger already initialized")
	ErrNotEnoughFunds                           = Error.Error("not enough funds")
	ErrInvalidTransactionSignature              = Error.Error("invalid transaction signature")
	ErrInvalidTransactionHash                   = Error.Error("invalid transaction hash")
	ErrTransactionAlreadyInLedger               = Error.Error("transaction already in ledger")
	ErrTransactionNotFound                      = Error.Error("previous not found")
	ErrPreviousTransactionNotFound              = Error.Error("previous transaction not found")
	ErrHeadTransactionNotFound                  = Error.Error("head transaction not found")
	ErrPreviousTransactionIsNotHead             = Error.Error("previous transaction is not the chain head")
	ErrSendTransactionIsNotPending              = Error.Error("send transaction is not pending")
	ErrOpenTransactionNotFound                  = Error.Error("open transaction not found")
	ErrAddressDoesNotMatchPubKey                = Error.Error("address does not match public key")
	ErrSendReceiveTransactionsNotLinked         = Error.Error("send and receive transaction not linked")
	ErrSendReceiveTransactionsCantBeSameAddress = Error.Error("send and receive can not be on the same address")
	ErrSentAmountDiffersFromReceivedAmount      = Error.Error("sent amount differs from received amount")
	ErrInvalidReceiveTransaction                = Error.Error("invalid receive transaction")
	ErrInvalidSendTransaction                   = Error.Error("invalid send transaction")
)

//go:generate mockgen -destination=../tests/mock_ledger.go -package=tests github.com/msaldanha/realChain/ledger Ledger

type Ledger interface {
	Initialize(genesisTx *transaction.Transaction) (error)
	GetLastTransaction(address string) (*transaction.Transaction, error)
	GetTransaction(hash string) (*transaction.Transaction, error)
	GetAddressStatement(address string) ([]*transaction.Transaction, error)
	Register(sendTx *transaction.Transaction, receiveTx *transaction.Transaction) (err error)
	VerifyTransaction(tx *transaction.Transaction, isNew bool) error
	VerifyTransactions(sendTx *transaction.Transaction, receiveTx *transaction.Transaction) (error)
}
