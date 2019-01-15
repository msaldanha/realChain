package wallet

import (
	"github.com/msaldanha/realChain/address"
	"github.com/msaldanha/realChain/keyvaluestore"
	"github.com/msaldanha/realChain/ledger"
	"github.com/msaldanha/realChain/errors"
)

const ErrAddressNotManagedByThisWallet       = errors.Error("address not managed by this wallet")

type Wallet struct {
	ld        ledger.Ledger
	ts        *ledger.TransactionStore
	addresses keyvaluestore.Storer
}

func New(txStore *ledger.TransactionStore, addressStore keyvaluestore.Storer, ld ledger.Ledger) (*Wallet) {
	return &Wallet{ld:ld, ts: txStore, addresses: addressStore}
}

func (wa *Wallet) CreateSendTransaction(from, to string, amount float64) (*ledger.Transaction, error) {
	fromTipTx, err := wa.ts.Retrieve(from)

	if err != nil {
		return nil, err
	}

	if fromTipTx == nil {
		return nil, ErrAddressNotManagedByThisWallet
	}

	if valid, err := address.IsValid(to); !valid {
		return nil, err
	}

	if fromTipTx.Balance - amount < 0 {
		return nil, ledger.ErrNotEnoughFunds
	}

	tx, err := wa.createSendTransaction(fromTipTx, []byte(to), amount)
	if err != nil {
		return nil, err
	}

	_, err = wa.ts.Store(tx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (wa *Wallet) GetAddressStatement(address string) ([]*ledger.Transaction, error) {
	return wa.ld.GetAddressStatement(address)
}

func (wa *Wallet) GetLastTransaction(address string) (*ledger.Transaction, error) {
	return wa.ld.GetLastTransaction(address)
}

func (wa *Wallet) GetAddresses() ([]*address.Address, error) {
	addrs, err := wa.addresses.GetAll()
	if err != nil {
		return nil, err
	}
	addresses := make([]*address.Address, 0)
	for _, v := range addrs {
		addresses = append(addresses, address.NewAddressFromBytes(v))
	}
	return addresses, nil
}

func (wa *Wallet) CreateAddress() (*address.Address, error) {
	addr, err := address.NewAddressWithKeys()
	if err != nil {
		return nil, err
	}

	wa.addresses.Put(addr.Address, addr.ToBytes())
	return addr, nil
}

func (wa *Wallet) createSendTransaction(fromTip *ledger.Transaction, to []byte, amount float64) (*ledger.Transaction, error) {
	addr, err := wa.GetAddress(fromTip.Address)
	if err != nil {
		return nil, err
	}

	send := ledger.NewSendTransaction()
	send.Address = fromTip.Address
	send.Link = to
	send.Previous = fromTip.Hash
	send.Balance = fromTip.Balance - amount
	send.PubKey = fromTip.PubKey

	send.SetPow()

	if err := send.Sign(addr.Keys.ToEcdsaPrivateKey()); err != nil {
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
