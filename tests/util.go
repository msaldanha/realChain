package tests

import (
	"github.com/msaldanha/realChain/ledger"
	. "github.com/onsi/gomega"
	"time"
	"github.com/msaldanha/realChain/keyvaluestore"
	"fmt"
	"strings"
	"github.com/msaldanha/realChain/address"
)

func AssertCommonVal(val ledger.Validator, tx *ledger.Transaction) {
	ok, err := val.IsFilled(tx)
	Expect(ok).To(BeFalse())
	Expect(err).NotTo(BeNil())
	Expect(err).To(Equal(ledger.ErrInvalidTransactionTimestamp))

	tx.Timestamp = time.Now().Unix()

	ok, err = val.IsFilled(tx)
	Expect(ok).To(BeFalse())
	Expect(err).NotTo(BeNil())
	Expect(err).To(Equal(ledger.ErrTransactionAddressCantBeEmpty))

	tx.Address = []byte("xxxxxxxxxxxxxxxxxxx")

	ok, err = val.IsFilled(tx)
	Expect(ok).To(BeFalse())
	Expect(err).NotTo(BeNil())
	Expect(err).To(Equal(ledger.ErrPreviousTransactionCantBeEmpty))

	tx.Previous = []byte("yyyyyyyyyyyyyyyyyyyy")

	ok, err = val.IsFilled(tx)
	Expect(ok).To(BeFalse())
	Expect(err).NotTo(BeNil())
	Expect(err).To(Equal(ledger.ErrTransactionSignatureCantBeEmpty))

	tx.Signature = []byte("ssssssssssssssssssssss")

	ok, err = val.IsFilled(tx)
	Expect(ok).To(BeFalse())
	Expect(err).NotTo(BeNil())
	Expect(err).To(Equal(ledger.ErrTransactionPowNonceCantBeZero))

	tx.PowNonce = 1

	ok, err = val.IsFilled(tx)
	Expect(ok).To(BeFalse())
	Expect(err).NotTo(BeNil())
	Expect(err).To(Equal(ledger.ErrTransactionHashCantBeEmpty))

	tx.SetHash()

	ok, err = val.IsFilled(tx)
	Expect(ok).To(BeFalse())
	Expect(err).NotTo(BeNil())
	Expect(err).To(Equal(ledger.ErrPubKeyCantBeEmpty))

	tx.PubKey = []byte("ssssssssssssssssssssss")
}

func CreateNonEmptyMemoryStore() *keyvaluestore.MemoryKeyValueStore {
	ms := keyvaluestore.NewMemoryKeyValueStore()
	tx := &ledger.Transaction{}
	ms.Put("genesis", tx.ToBytes())
	return ms
}

func DumpTxChain(txChain []*ledger.Transaction) {
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

func CreateGenesisTransaction(balance float64) (*ledger.Transaction, *address.Address) {
	genesisTx := ledger.NewOpenTransaction()
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


func SendFunds(ld ledger.Ledger, sendAddr *address.Address, prevSendTx *ledger.Transaction, prevReceiveTx *ledger.Transaction,
		receiveAddr *address.Address, amount float64) (*ledger.Transaction, *ledger.Transaction) {
	sendTx, err := ledger.CreateSendTransaction(prevSendTx, sendAddr, receiveAddr.Address, amount)
	Expect(err).To(BeNil())

	receiveTx, err := ledger.CreateReceiveTransaction(sendTx, amount, receiveAddr, prevReceiveTx)
	Expect(err).To(BeNil())

	err = ld.Register(sendTx, receiveTx)
	Expect(err).To(BeNil())

	return sendTx, receiveTx
}