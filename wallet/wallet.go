package wallet

import (
	"context"
	"github.com/msaldanha/realChain/address"
	"github.com/msaldanha/realChain/errors"
	"github.com/msaldanha/realChain/keyvaluestore"
	"github.com/msaldanha/realChain/ledger"
	"google.golang.org/grpc"
)

const (
	ErrAddressNotManagedByThisWallet = errors.Error("address not managed by this wallet")
)

type Wallet struct {
	ld        ledger.LedgerClient
	addresses keyvaluestore.Storer
	ctx       context.Context
	opts      grpc.CallOption
}

func New(addressStore keyvaluestore.Storer, ld ledger.LedgerClient) *Wallet {
	return &Wallet{ld: ld, addresses: addressStore, ctx: context.Background(), opts: &grpc.EmptyCallOption{}}
}

func (wa *Wallet) Transfer(from, to string, amount float64) (*ledger.Transaction, error) {
	if valid, err := address.IsValid(to); !valid {
		return nil, err
	}

	fromAddr, err := wa.getAddress(from)
	if err != nil {
		return nil, err
	}

	toAddr, err := wa.getAddress(to)
	if err != nil {
		return nil, err
	}

	toTipTx, err := wa.GetLastTransaction(to)
	if err != nil {
		return nil, err
	}

	fromTipTx, err := wa.GetLastTransaction(from)
	if err != nil {
		return nil, err
	}
	if fromTipTx == nil {
		return nil, ledger.ErrNotEnoughFunds
	}

	if fromTipTx.Balance-amount < 0 {
		return nil, ledger.ErrNotEnoughFunds
	}

	sendTx, err := ledger.CreateSendTransaction(fromTipTx, fromAddr, toAddr.Address, amount)
	if err != nil {
		return nil, err
	}

	receiveTx, err := ledger.CreateReceiveTransaction(sendTx, amount, toAddr, toTipTx)
	if err != nil {
		return nil, err
	}

	_, err = wa.ld.Register(wa.ctx, &ledger.RegisterRequest{SendTx: sendTx, ReceiveTx: receiveTx}, wa.opts)
	if err != nil {
		return nil, err
	}

	return receiveTx, nil
}

func (wa *Wallet) GetAddressStatement(addr string) ([]*ledger.Transaction, error) {
	result, err := wa.ld.GetAddressStatement(wa.ctx, &ledger.GetAddressStatementRequest{Address: addr}, wa.opts)
	if err != nil {
		return nil, err
	}

	return result.Txs, nil
}

func (wa *Wallet) GetLastTransaction(addr string) (*ledger.Transaction, error) {
	result, err := wa.ld.GetLastTransaction(wa.ctx, &ledger.GetLastTransactionRequest{Address: addr}, wa.opts)
	if err != nil {
		return nil, err
	}
	return result.Tx, nil
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

	err = wa.addresses.Put(addr.Address, addr.ToBytes())
	if err != nil {
		return nil, err
	}

	return addr, nil
}

func (wa *Wallet) getAddress(addr string) (*address.Address, error) {
	addrBytes, ok, err := wa.addresses.Get(addr)
	if ok == false && err == nil {
		return nil, ErrAddressNotManagedByThisWallet
	}
	if err != nil {
		return nil, err
	}

	return address.NewAddressFromBytes(addrBytes), nil
}

func (wa *Wallet) GetAddress(addressBytes []byte) (*address.Address, error) {
	addr, ok, err := wa.addresses.Get(string(addressBytes))
	if !ok {
		return nil, err
	}
	return address.NewAddressFromBytes(addr), nil
}
