package ledger

import (
	"github.com/msaldanha/realChain/keyvaluestore"
)

type TransactionStore struct {
	store            keyvaluestore.Storer
	validatorCreator ValidatorCreator
}

func NewTransactionStore(store keyvaluestore.Storer, validatorCreator ValidatorCreator) (*TransactionStore) {
	a := &TransactionStore{store: store, validatorCreator: validatorCreator}
	return a
}

func (ts *TransactionStore) isValid(tx *Transaction) (bool, error) {
	val := ts.validatorCreator.CreateValidatorForTransaction(tx.Type, ts.store)
	return val.IsValid(tx)
}

func (ts *TransactionStore) Store(tx *Transaction) (*Transaction, error) {
	if ok, err := ts.isValid(tx); !ok {
		return nil, err
	}

	err := ts.store.Put(string(tx.Hash), tx.ToBytes())
	if err != nil {
		return nil, err
	}

	err = ts.store.Put(string(tx.Address), tx.ToBytes())
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (ts *TransactionStore) Retrieve(hash string) (*Transaction, error) {
	value, _, err := ts.GetTransaction(hash)
	if err != nil {
		return nil, err
	}
	return value, nil
}


func (ts *TransactionStore) GetTransactionChain(txHash string, includeAll bool) ([]*Transaction, error) {
	tx, ok, _ := ts.GetTransaction(txHash)
	chain := make([]*Transaction, 0)
	for ok {
		chain = append(chain[:0], append([]*Transaction{tx}, chain[0:]...)...)
		if len(tx.Previous) > 0 {
			tx, ok, _ = ts.GetTransaction(string(tx.Previous))
		} else if tx.Type == Transaction_OPEN && len(tx.Link) > 0 && includeAll {
			tx, ok, _ = ts.GetTransaction(string(tx.Link))
		} else {
			break
		}
	}
	return chain, nil
}

func (ts *TransactionStore) GetTransaction(txHash string) (*Transaction, bool, error) {
	tx, ok, err := ts.store.Get(txHash)
	if tx == nil {
		return nil, ok, err
	}
	return NewTransactionFromBytes(tx), ok, err
}

func (ts *TransactionStore) IsEmpty() (bool) {
	return ts.store.IsEmpty()
}
