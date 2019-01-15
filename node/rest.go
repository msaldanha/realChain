package node

import (
	"github.com/msaldanha/realChain/errors"
	"github.com/kataras/iris"
	"github.com/msaldanha/realChain/address"
	log "github.com/sirupsen/logrus"
	"encoding/hex"
	"github.com/msaldanha/realChain/ledger"
	"strings"
)

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

func mapToAddressDto(acc *address.Address, balance float64) *AddressDto {
	addrDto := &AddressDto{Address: acc.Address, Keys: &KeyPairDto{}}
	addrDto.Keys.PublicKey = hex.EncodeToString(acc.Keys.PublicKey)
	addrDto.Keys.PrivateKey = hex.EncodeToString(acc.Keys.PrivateKey)
	addrDto.Balance = balance
	return addrDto
}

func mapToTransactionDtos(txchain []*ledger.Transaction) []*TransactionDto {
	txs := make([]*TransactionDto, 0)
	for _, v := range txchain {
		tx := mapToTransactionDto(v)
		txs = append(txs, tx)
	}
	return txs
}

func mapToTransactionDto(tx *ledger.Transaction) *TransactionDto {
	return &TransactionDto{
		Id:        string(tx.Hash),
		Address:   string(tx.Address),
		Type:      tx.Type.String(),
		Previous:  string(tx.Previous),
		Balance:   tx.Balance,
		Link:      string(tx.Link),
		Timestamp: tx.Timestamp,
	}
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
		if _, ok := err.(errors.Error); ok {
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
