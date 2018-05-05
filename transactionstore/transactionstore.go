package transactionstore

import (
	"github.com/msaldanha/realChain/transaction"
	"github.com/msaldanha/realChain/keyvaluestore"
)

type TransactionStore struct {
	store            keyvaluestore.Storer
	validatorCreator transaction.ValidatorCreator
}

func New(store keyvaluestore.Storer, validatorCreator transaction.ValidatorCreator) (*TransactionStore) {
	a := &TransactionStore{store: store, validatorCreator: validatorCreator}
	return a
}

func (ts *TransactionStore) isValid(tx *transaction.Transaction) (bool, error) {
	if !tx.Type.IsValid(){
		return false, transaction.ErrInvalidTransactionType
	}
	val := ts.validatorCreator.CreateValidatorForTransaction(tx.Type, ts.store)
	return val.IsValid(tx)
}

func (ts *TransactionStore) Store(tx *transaction.Transaction) (*transaction.Transaction, error) {
	if ok, err := ts.isValid(tx); !ok {
		return nil, err
	}
	ts.store.Put(string(tx.Hash), tx.ToBytes())
	ts.store.Put(string(tx.Address), tx.ToBytes())
	return tx, nil
}

func (ts *TransactionStore) Retrieve(hash string) (*transaction.Transaction, error) {
	value, _, err := ts.GetTransaction(hash)
	if err != nil {
		return nil, err
	}
	return value, nil
}


func (ts *TransactionStore) GetTransactionChain(txHash string, includeAll bool) ([]*transaction.Transaction, error) {
	tx, ok, _ := ts.GetTransaction(txHash)
	chain := []*transaction.Transaction{}
	for ok {
		chain = append(chain[:0], append([]*transaction.Transaction{tx}, chain[0:]...)...)
		if len(tx.Previous) > 0 {
			tx, ok, _ = ts.GetTransaction(string(tx.Previous))
		} else if tx.Type == transaction.OPEN && len(tx.Link) > 0 && includeAll {
			tx, ok, _ = ts.GetTransaction(string(tx.Link))
		} else {
			break
		}
	}
	return chain, nil
}

func (ts *TransactionStore) GetTransaction(txHash string) (*transaction.Transaction, bool, error) {
	tx, ok, err := ts.store.Get(txHash)
	if tx == nil {
		return nil, ok, err
	}
	return transaction.NewTransactionFromBytes(tx), ok, err
}

func (ts *TransactionStore) IsEmpty() (bool) {
	return ts.store.IsEmpty()
}
