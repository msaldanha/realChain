package wallet

import (
	"github.com/msaldanha/realChain/transaction"
	"github.com/msaldanha/realChain/address"
	"github.com/msaldanha/realChain/transactionstore"
	"github.com/msaldanha/realChain/keyvaluestore"
)

type Wallet struct {
	ts        *transactionstore.TransactionStore
	addresses keyvaluestore.Storer
}

func New(txStore *transactionstore.TransactionStore, addressStore keyvaluestore.Storer) (*Wallet) {
	return &Wallet{ts:txStore, addresses:addressStore}
}

func (wa *Wallet) CreateSendTransaction(from, to string, amount float64) (*transaction.Transaction, error) {
	fromTipTx, err := wa.ts.Retrieve(from)

	if err != nil {
		return nil, err
	}

	addr := address.New()
	addr.Address = to
	if valid, err := addr.IsValid(); !valid {
		return nil, err
	}

	tx, err := wa.createSendTransaction(fromTipTx, []byte(to), amount)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (wa *Wallet) createSendTransaction(fromTip *transaction.Transaction, to []byte, amount float64) (*transaction.Transaction, error) {
	addr, err := wa.GetAddress(fromTip.Address)
	if err != nil {
		return nil, err
	}

	send := transaction.NewSendTransaction()
	send.Address = fromTip.Address
	send.Link = to
	send.Previous = fromTip.Hash
	send.Balance = fromTip.Balance - amount
	send.PubKey = fromTip.PubKey

	send.SetPow()

	if err := send.Sign(addr.Keys.ToEcdsaPrivateKey()); err != nil {
		return nil, err
	}

	_, err = wa.ts.Store(send)
	if err != nil {
		return nil, err
	}

	return send, nil
}

func (wa *Wallet) GetAddress(addressBytes []byte) (*address.Address, error) {
	addr, ok, err := wa.addresses.Get(string(addressBytes))
	if !ok {
		return nil, err
	}
	return address.NewAddressFromBytes(addr), nil
}