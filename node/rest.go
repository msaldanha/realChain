package node

import (
	"github.com/kataras/iris"
	"github.com/msaldanha/realChain/ledger"
	"strings"
	"encoding/hex"
	"github.com/msaldanha/realChain/keypair"
	"github.com/msaldanha/realChain/transaction"
	log "github.com/sirupsen/logrus"
	"github.com/msaldanha/realChain/Error"
)

type RestServer struct {
	iris *iris.Application
	ld   *ledger.Ledger
	url  string
}

type KeyPairDto struct {
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
}

type AddressDto struct {
	Address string      `json:"address"`
	Keys    *KeyPairDto `json:"keys"`
	Balance float64     `json:"balance"`
}

type SendDto struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
	TxId   string  `json:"txid"`
}

type ErrorDto struct {
	Message string `json:"message"`
}

type TransactionDto struct {
	Id        string  `json:"id"`
	Address   string  `json:"address"`
	Type      string  `json:"type"`
	Previous  string  `json:"previous"`
	Balance   float64 `json:"balance"`
	Link      string  `json:"link"`
	Timestamp int64   `json:"timestamp"`
}

func NewRestServer(l *ledger.Ledger, url string) *RestServer {
	irisApp := iris.New()
	return &RestServer{iris: irisApp, ld: l, url: url}
}

func (rs *RestServer) Run() error {
	log.Info("Rest server starting")
	rs.iris.Get("/wallet/addresses", rs.getAddresses())
	rs.iris.Post("/wallet/addresses", rs.createAddresses())
	rs.iris.Get("/wallet/addresses/{address:string}", rs.getAddressByAddress())
	rs.iris.Get("/wallet/addresses/{address:string}/statement", rs.getAddressesStatement())
	rs.iris.Post("/wallet/tx", rs.sendFunds())
	return rs.iris.Run(iris.Addr(rs.url), iris.WithoutServerError(iris.ErrServerClosed))
}

func (rs *RestServer) getAddressesStatement() iris.Handler {
	return func(ctx iris.Context) {
		logRequest(ctx)
		addr := ctx.Params().Get("address")
		txs, err := rs.ld.GetAddressStatement(addr)

		if hasError(ctx, err) {
			return
		}

		ctx.JSON(mapToTransactionDtos(txs))
	}
}

func (rs *RestServer) getAddressByAddress() iris.Handler {
	return func(ctx iris.Context) {
		logRequest(ctx)
		addr := ctx.Params().Get("address")
		acc, err := rs.ld.GetAddress([]byte(addr))

		if hasError(ctx, err) {
			return
		}

		if acc == nil {
			acc = &ledger.Address{}
			acc.Keys = &keypair.KeyPair{}
		}

		tx, err := rs.ld.GetLastTransaction(addr)
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

func (rs *RestServer) getAddresses() iris.Handler {
	return func(ctx iris.Context) {
		logRequest(ctx)
		acc, err := rs.ld.GetAddresses()

		if hasError(ctx, err) {
			return
		}

		addresses := make([]*AddressDto, 0)
		for _, v := range acc {
			var balance float64 = 0
			tx, err := rs.ld.GetLastTransaction(v.Address)
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

func (rs *RestServer) createAddresses() iris.Handler {
	return func(ctx iris.Context) {
		logRequest(ctx)
		acc, err := rs.ld.CreateAddress()

		if hasError(ctx, err) {
			return
		}
		ctx.JSON(mapToAddressDto(acc, 0))
	}
}

func (rs *RestServer) sendFunds() iris.Handler {
	return func(ctx iris.Context) {
		logRequest(ctx)
		send := &SendDto{}
		ctx.ReadJSON(send)
		tx, err := rs.ld.CreateSendTransaction(send.From, send.To, send.Amount)
		if hasError(ctx, err) {
			return
		}
		err = rs.ld.HandleTransaction(tx)
		if hasError(ctx, err) {
			return
		}
		send.TxId = string(tx.Hash)
		ctx.JSON(send)
	}
}

func mapToAddressDto(acc *ledger.Address, balance float64) *AddressDto {
	addrDto := &AddressDto{Address: acc.Address, Keys: &KeyPairDto{}}
	addrDto.Keys.PublicKey = hex.EncodeToString(acc.Keys.PublicKey)
	addrDto.Keys.PrivateKey = hex.EncodeToString(acc.Keys.PrivateKey)
	addrDto.Balance = balance
	return addrDto
}

func mapToTransactionDtos(txchain []*transaction.Transaction) []*TransactionDto {
	txs := make([]*TransactionDto, 0)
	for _, v := range txchain {
		tx := &TransactionDto{
			Id:        string(v.Hash),
			Address:   string(v.Address),
			Type:      v.Type.String(),
			Previous:  string(v.Previous),
			Balance:   v.Balance,
			Link:      string(v.Link),
			Timestamp: v.Timestamp,
		}
		txs = append(txs, tx)
	}
	return txs
}

func hasError(ctx iris.Context, err error) bool {
	if err != nil {
		setError(ctx, err)
		return true
	}
	return false
}

func setError(ctx iris.Context, err error) {
	if strings.Contains(err.Error(), "not found") {
		ctx.StatusCode(404)
	} else {
		if _, ok := err.(Error.Error); ok {
			ctx.StatusCode(400)
		} else {
			ctx.StatusCode(500)
		}
	}
	log.Error(err)
	ctx.JSON(errorFor(err))
}

func errorFor(err error) *ErrorDto {
	return &ErrorDto{Message: err.Error()}
}

func logRequest(ctx iris.Context) {
	log.Infof("Rest request for %s %s", ctx.Method(), ctx.Path())
}
