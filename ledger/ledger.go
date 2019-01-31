package ledger

import (
	"github.com/msaldanha/realChain/errors"
)

const (
	ErrLedgerAlreadyInitialized                 = errors.Error("ledger already initialized")
	ErrNotEnoughFunds                           = errors.Error("not enough funds")
	ErrInvalidTransactionSignature              = errors.Error("invalid transaction signature")
	ErrInvalidTransactionHash                   = errors.Error("invalid transaction hash")
	ErrTransactionAlreadyInLedger               = errors.Error("transaction already in ledger")
	ErrTransactionNotFound                      = errors.Error("previous not found")
	ErrPreviousTransactionNotFound              = errors.Error("previous transaction not found")
	ErrHeadTransactionNotFound                  = errors.Error("head transaction not found")
	ErrPreviousTransactionIsNotHead             = errors.Error("previous transaction is not the chain head")
	ErrSendTransactionIsNotPending              = errors.Error("send transaction is not pending")
	ErrOpenTransactionNotFound                  = errors.Error("open transaction not found")
	ErrAddressDoesNotMatchPubKey                = errors.Error("address does not match public key")
	ErrSendReceiveTransactionsNotLinked         = errors.Error("send and receive transaction not linked")
	ErrSendReceiveTransactionsCantBeSameAddress = errors.Error("send and receive can not be on the same address")
	ErrSentAmountDiffersFromReceivedAmount      = errors.Error("sent amount differs from received amount")
	ErrInvalidReceiveTransaction                = errors.Error("invalid receive transaction")
	ErrInvalidSendTransaction                   = errors.Error("invalid send transaction")
)

//go:generate protoc -I.. ledger/ledgerserver.proto --go_out=plugins=grpc:../
//go:generate mockgen -destination=../tests/mock_ledger.go -package=tests github.com/msaldanha/realChain/ledger Ledger
//go:generate mockgen -destination=../tests/mock_ledgerclient.go -package=tests github.com/msaldanha/realChain/ledger LedgerClient

type Ledger interface {
	Initialize(genesisTx *Transaction) error
	GetLastTransaction(address string) (*Transaction, error)
	GetTransaction(hash string) (*Transaction, error)
	GetAddressStatement(address string) ([]*Transaction, error)
	Register(sendTx *Transaction, receiveTx *Transaction) error
	VerifyTransaction(tx *Transaction, isNew bool) error
	Verify(sendTx *Transaction, receiveTx *Transaction) error
}
