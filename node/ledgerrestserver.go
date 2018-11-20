package node

import (
	"github.com/kataras/iris"
	log "github.com/sirupsen/logrus"
	"net"
	"github.com/msaldanha/realChain/ledger"
)

type LedgerRestServer struct {
	iris *iris.Application
	ld ledger.Ledger
	url  string
	conn *net.UDPConn
}

func NewLedgerRestServer(ld ledger.Ledger, url string) (*LedgerRestServer, error) {
	irisApp := iris.New()
	return &LedgerRestServer{iris: irisApp, ld: ld, url: url}, nil
}

func (rs *LedgerRestServer) Run() error {
	log.Info("Ledger rest server starting")
	rs.iris.Get("/ledger/addresses/{address:string}/tx/last", rs.getLastTransaction())
	rs.iris.Get("/ledger/tx/{hash:string}", rs.getTransaction())
	rs.iris.Get("/ledger/addresses/{address:string}/statement", rs.getAddressesStatement())
	return rs.iris.Run(iris.Addr(rs.url), iris.WithoutServerError(iris.ErrServerClosed))
}

func (rs *LedgerRestServer) getAddressesStatement() iris.Handler {
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

func (rs *LedgerRestServer) getLastTransaction() iris.Handler {
	return func(ctx iris.Context) {
		logRequest(ctx)
		addr := ctx.Params().Get("address")
		tx, err := rs.ld.GetLastTransaction(addr)

		if hasError(ctx, err) {
			return
		}

		ctx.JSON(mapToTransactionDto(tx))
	}
}

func (rs *LedgerRestServer) getTransaction() iris.Handler {
	return func(ctx iris.Context) {
		logRequest(ctx)
		hash := ctx.Params().Get("hash")
		tx, err := rs.ld.GetTransaction(hash)

		if hasError(ctx, err) {
			return
		}

		ctx.JSON(mapToTransactionDto(tx))
	}
}
