package node

import (
	"github.com/kataras/iris"
	"github.com/msaldanha/realChain/keypair"
	log "github.com/sirupsen/logrus"
	"net"
	"github.com/msaldanha/realChain/wallet"
	"github.com/msaldanha/realChain/address"
)

type WalletRestServer struct {
	iris *iris.Application
	wa   *wallet.Wallet
	url  string
	conn *net.UDPConn
}

func NewWalletRestServer(wa *wallet.Wallet, url string) (*WalletRestServer, error) {
	irisApp := iris.New()
	return &WalletRestServer{iris: irisApp, wa: wa, url: url}, nil
}

func (rs *WalletRestServer) Run() error {
	log.Info("Wallet rest server starting")
	rs.iris.Get("/wallet/addresses", rs.getAddresses())
	rs.iris.Post("/wallet/addresses", rs.createAddresses())
	rs.iris.Get("/wallet/addresses/{address:string}", rs.getAddressByAddress())
	rs.iris.Get("/wallet/addresses/{address:string}/statement", rs.getAddressesStatement())
	rs.iris.Post("/wallet/tx", rs.sendFunds())
	return rs.iris.Run(iris.Addr(rs.url), iris.WithoutServerError(iris.ErrServerClosed))
}

func (rs *WalletRestServer) getAddressesStatement() iris.Handler {
	return func(ctx iris.Context) {
		logRequest(ctx)
		addr := ctx.Params().Get("address")
		txs, err := rs.wa.GetAddressStatement(addr)

		if hasError(ctx, err) {
			return
		}

		ctx.JSON(mapToTransactionDtos(txs))
	}
}

func (rs *WalletRestServer) getAddressByAddress() iris.Handler {
	return func(ctx iris.Context) {
		logRequest(ctx)
		addr := ctx.Params().Get("address")
		acc, err := rs.wa.GetAddress([]byte(addr))

		if hasError(ctx, err) {
			return
		}

		if acc == nil {
			acc = &address.Address{}
			acc.Keys = &keypair.KeyPair{}
		}

		tx, err := rs.wa.GetLastTransaction(addr)
		if hasError(ctx, err) {
			return
		}

		if tx == nil {
			ctx.StatusCode(404)
			return
		}

		acc.Address = string(tx.Address)
		acc.Keys.PublicKey = tx.PubKey

		ctx.JSON(mapToAddressDto(acc, tx.Balance))
	}
}

func (rs *WalletRestServer) getAddresses() iris.Handler {
	return func(ctx iris.Context) {
		logRequest(ctx)
		acc, err := rs.wa.GetAddresses()

		if hasError(ctx, err) {
			return
		}

		addresses := make([]*AddressDto, 0)
		for _, v := range acc {
			var balance float64 = 0
			tx, err := rs.wa.GetLastTransaction(v.Address)
			if hasError(ctx, err) {
				return
			}
			if tx != nil {
				balance = tx.Balance
			}
			addresses = append(addresses, mapToAddressDto(v, balance))
		}

		ctx.JSON(addresses)
	}
}

func (rs *WalletRestServer) createAddresses() iris.Handler {
	return func(ctx iris.Context) {
		logRequest(ctx)
		acc, err := rs.wa.CreateAddress()

		if hasError(ctx, err) {
			return
		}
		ctx.JSON(mapToAddressDto(acc, 0))
	}
}

func (rs *WalletRestServer) sendFunds() iris.Handler {
	return func(ctx iris.Context) {
		logRequest(ctx)
		send := &SendDto{}
		ctx.ReadJSON(send)
		tx, err := rs.wa.SendFunds(send.From, send.To, send.Amount)
		if hasError(ctx, err) {
			return
		}

		send.TxId = string(tx.Hash)
		ctx.JSON(send)
	}
}
