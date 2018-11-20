package ledger

import (
	"github.com/msaldanha/realChain/transactionstore"
	"github.com/msaldanha/realChain/transaction"
	"github.com/msaldanha/realChain/address"
	"github.com/msaldanha/realChain/keyvaluestore"
)

type RemoteLedger struct {
	ts        *transactionstore.TransactionStore
	addresses keyvaluestore.Storer
}

func (ld *RemoteLedger) Initialize(initialBalance float64) (*transaction.Transaction, *address.Address, error) {
	return nil, nil, ErrInvalidOperation
}

func (ld *RemoteLedger) Receive(send *transaction.Transaction) (string, error) {
	return "", ErrInvalidOperation
}

func (ld *RemoteLedger) HandleTransactionBytes(txBytes []byte) (*transaction.Transaction, error) {
	return nil, ErrInvalidOperation
}

func (ld *RemoteLedger) HandleTransaction(tx *transaction.Transaction) (ret *transaction.Transaction, err error) {
	return nil, ErrInvalidOperation
}

func (ld *RemoteLedger) GetLastTransaction(address string) (*transaction.Transaction, error) {
	return nil, ErrInvalidOperation
}

func (ld *RemoteLedger) GetTransaction(hash string) (*transaction.Transaction, error) {
	return nil, ErrInvalidOperation
}

func (ld *RemoteLedger) GetAddressStatement(address string) ([]*transaction.Transaction, error) {
	return nil, ErrInvalidOperation
}

func (ld *RemoteLedger) AddAddress(addr *address.Address) error {
	return ErrInvalidOperation
}

func (ld *RemoteLedger) VerifyTransaction(tx *transaction.Transaction, isNew bool) (error) {
	return ErrInvalidOperation
}
