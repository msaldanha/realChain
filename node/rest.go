package node

import (
	"github.com/kataras/iris"
	"github.com/msaldanha/realChain/ledger"
	"strings"
	"encoding/hex"
	"github.com/msaldanha/realChain/keypair"
	"github.com/msaldanha/realChain/block"
	"log"
	"github.com/msaldanha/realChain/Error"
)

type RestServer struct {
	iris *iris.Application
	ld   *ledger.Ledger
}

type KeyPairDto struct {
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
}

type AccountDto struct {
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
	Type      string  `json:"type"`
	Previous  string  `json:"previous"`
	Balance   float64 `json:"balance"`
	Link      string  `json:"link"`
	Timestamp int64   `json:"timestamp"`
}

func NewRestServer(l *ledger.Ledger) *RestServer {
	irisApp := iris.New()

	return &RestServer{iris: irisApp, ld: l}
}

func (rs *RestServer) Run() error {
	rs.iris.Get("/wallet/accounts", rs.getAccounts())
	rs.iris.Post("/wallet/accounts", rs.createAccount())
	rs.iris.Get("/wallet/accounts/{address:string}", rs.getAccountByAddress())
	rs.iris.Get("/wallet/accounts/{address:string}/statement", rs.getAccountStatementByAddress())
	rs.iris.Post("/wallet/tx", rs.sendFunds())
	return rs.iris.Run(iris.Addr(":1300"), iris.WithoutServerError(iris.ErrServerClosed))
}

func (rs *RestServer) getAccountStatementByAddress() iris.Handler {
	return func(ctx iris.Context) {
		addr := ctx.Params().Get("address")
		acc, err := rs.ld.GetAccountStatement(addr)

		if hasError(ctx, err) {
			return
		}

		ctx.JSON(mapToTransactionDtos(acc))
	}
}

func (rs *RestServer) getAccountByAddress() iris.Handler {
	return func(ctx iris.Context) {
		addr := ctx.Params().Get("address")
		acc, err := rs.ld.GetAccount([]byte(addr))

		if hasError(ctx, err) {
			return
		}

		if acc == nil {
			acc = &ledger.Account{}
			acc.Keys = &keypair.KeyPair{}
		}

		blk, err := rs.ld.GetLastTransaction(addr)
		if hasError(ctx, err) {
			return
		}

		if blk == nil {
			ctx.StatusCode(404)
			return
		}

		acc.Address = string(blk.Account)
		acc.Keys.PublicKey = blk.PubKey

		ctx.JSON(mapToAccountDto(acc, blk.Balance))
	}
}

func (rs *RestServer) getAccounts() iris.Handler {
	return func(ctx iris.Context) {
		acc, err := rs.ld.GetAccounts()

		if hasError(ctx, err) {
			return
		}

		accounts := make([]*AccountDto, 0)
		for _, v := range acc {
			var balance float64 = 0
			blk, err := rs.ld.GetLastTransaction(v.Address)
			if hasError(ctx, err) {
				return
			}
			if blk != nil {
				balance = blk.Balance
			}
			accounts = append(accounts, mapToAccountDto(v, balance))
		}

		ctx.JSON(accounts)
	}
}

func (rs *RestServer) createAccount() iris.Handler {
	return func(ctx iris.Context) {
		acc, err := rs.ld.CreateAccount()

		if hasError(ctx, err) {
			return
		}
		ctx.JSON(mapToAccountDto(acc, 0))
	}
}

func (rs *RestServer) sendFunds() iris.Handler {
	return func(ctx iris.Context) {
		send := &SendDto{}
		ctx.ReadJSON(send)
		id, err := rs.ld.Send(send.From, send.To, send.Amount)
		if hasError(ctx, err) {
			return
		}
		send.TxId = id
		ctx.JSON(send)
	}
}

func mapToAccountDto(acc *ledger.Account, balance float64) *AccountDto {
	accDto := &AccountDto{Address: acc.Address, Keys: &KeyPairDto{}}
	accDto.Keys.PublicKey = hex.EncodeToString(acc.Keys.PublicKey)
	accDto.Keys.PrivateKey = hex.EncodeToString(acc.Keys.PrivateKey)
	accDto.Balance = balance
	return accDto
}

func mapToTransactionDtos(blockchain []*block.Block) []*TransactionDto {
	txs := make([]*TransactionDto, 0)
	for _, v := range blockchain {
		tx := &TransactionDto{
			Id:        string(v.Hash),
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
	log.Println(err)
	ctx.JSON(errorFor(err))
}

func errorFor(err error) *ErrorDto {
	return &ErrorDto{Message: err.Error()}
}
