package tests

import (
	"github.com/msaldanha/realChain/ledger"
	. "github.com/onsi/gomega"
	"github.com/msaldanha/realChain/transaction"
	"time"
	"github.com/msaldanha/realChain/keyvaluestore"
	"fmt"
	"strings"
	"github.com/msaldanha/realChain/address"
)

func AssertCommonVal(val transaction.Validator, tx *transaction.Transaction) {
	ok, err := val.IsFilled(tx)
	Expect(ok).To(BeFalse())
	Expect(err).NotTo(BeNil())
	Expect(err).To(Equal(transaction.ErrInvalidTransactionTimestamp))

	tx.Timestamp = time.Now().Unix()

	ok, err = val.IsFilled(tx)
	Expect(ok).To(BeFalse())
	Expect(err).NotTo(BeNil())
	Expect(err).To(Equal(transaction.ErrTransactionAddressCantBeEmpty))

	tx.Address = []byte("xxxxxxxxxxxxxxxxxxx")

	ok, err = val.IsFilled(tx)
	Expect(ok).To(BeFalse())
	Expect(err).NotTo(BeNil())
	Expect(err).To(Equal(transaction.ErrPreviousTransactionCantBeEmpty))

	tx.Previous = []byte("yyyyyyyyyyyyyyyyyyyy")

	ok, err = val.IsFilled(tx)
	Expect(ok).To(BeFalse())
	Expect(err).NotTo(BeNil())
	Expect(err).To(Equal(transaction.ErrTransactionSignatureCantBeEmpty))

	tx.Signature = []byte("ssssssssssssssssssssss")

	ok, err = val.IsFilled(tx)
	Expect(ok).To(BeFalse())
	Expect(err).NotTo(BeNil())
	Expect(err).To(Equal(transaction.ErrTransactionPowNonceCantBeZero))

	tx.PowNonce = 1

	ok, err = val.IsFilled(tx)
	Expect(ok).To(BeFalse())
	Expect(err).NotTo(BeNil())
	Expect(err).To(Equal(transaction.ErrTransactionHashCantBeEmpty))

	tx.SetHash()

	ok, err = val.IsFilled(tx)
	Expect(ok).To(BeFalse())
	Expect(err).NotTo(BeNil())
	Expect(err).To(Equal(transaction.ErrPubKeyCantBeEmpty))

	tx.PubKey = []byte("ssssssssssssssssssssss")
}

func CreateNonEmptyMemoryStore() *keyvaluestore.MemoryKeyValueStore {
	ms := keyvaluestore.NewMemoryKeyValueStore()
	tx := &transaction.Transaction{}
	ms.Put("genesis", tx.ToBytes())
	return ms
}

func DumpTxChain(txChain []*transaction.Transaction) {
	fmt.Println("============= Transaction Chain Dump start =================")
	level := 0
	for _, v := range txChain {
		if len(v.Previous) == 0 {
			level = level + 1
		}
		fmt.Printf("%s %s %s %s %s %f\n", strings.Repeat("  ", level), v.Type, string(v.Address), string(v.Hash), string(v.Previous), v.Balance)
	}
	fmt.Println("============= End =================")
}

func CreateTestAddress() *address.Address {
	addr, _ := address.NewAddressWithKeys()
	return addr
}

func CreateSendTransaction(fromTip *transaction.Transaction, fromAddr *address.Address, to string, amount float64) (*transaction.Transaction, error) {
	send := transaction.NewSendTransaction()
	send.Address = fromTip.Address
	send.Link = []byte(to)
	send.Previous = fromTip.Hash
	send.Balance = fromTip.Balance - amount
	send.PubKey = fromTip.PubKey
	send.SetPow()
	send.Sign(fromAddr.Keys.ToEcdsaPrivateKey())
	return send, nil
}

func CreateReceiveTransaction(send *transaction.Transaction, amount float64, receiveAddr *address.Address,
	receiveTip *transaction.Transaction) (*transaction.Transaction, error) {

	var receive *transaction.Transaction
	if receiveTip != nil {
		receive = transaction.NewReceiveTransaction()
		receive.Previous = receiveTip.Hash
		receive.Balance = receiveTip.Balance + amount
		receive.Representative = receiveTip.Representative
		receive.PubKey = receiveTip.PubKey
	} else {
		receive = transaction.NewOpenTransaction()
		receive.Balance = amount
		receive.Representative = send.Link
		receive.PubKey = receiveAddr.Keys.PublicKey
	}

	receive.Address = send.Link
	receive.Link = send.Hash

	if err := receive.SetPow(); err != nil {
		return nil, err
	}

	if err := receive.Sign(receiveAddr.Keys.ToEcdsaPrivateKey()); err != nil {
		return nil, err
	}

	return receive, nil
}

func CreateGenesisTransaction(balance float64) (*transaction.Transaction, *address.Address) {
	genesisTx := transaction.NewOpenTransaction()
	addr, err := address.NewAddressWithKeys()
	Expect(err).To(BeNil())

	genesisTx.Address = []byte(addr.Address)
	genesisTx.Representative = genesisTx.Address
	genesisTx.Balance = balance
	genesisTx.PubKey = addr.Keys.PublicKey

	err = genesisTx.SetPow()
	Expect(err).To(BeNil())

	err = genesisTx.Sign(addr.Keys.ToEcdsaPrivateKey())
	Expect(err).To(BeNil())

	return genesisTx, addr
}


func SendFunds(ld ledger.Ledger, sendAddr *address.Address, prevSendTx *transaction.Transaction, prevReceiveTx *transaction.Transaction,
		receiveAddr *address.Address, amount float64) (*transaction.Transaction, *transaction.Transaction) {
	sendTx, err := CreateSendTransaction(prevSendTx, sendAddr, receiveAddr.Address, amount)
	Expect(err).To(BeNil())

	receiveTx, err := CreateReceiveTransaction(sendTx, amount, receiveAddr, prevReceiveTx)
	Expect(err).To(BeNil())

	err = ld.Register(sendTx, receiveTx)
	Expect(err).To(BeNil())

	return sendTx, receiveTx
}