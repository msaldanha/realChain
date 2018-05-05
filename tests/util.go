package tests

import (
	. "github.com/onsi/gomega"
	"github.com/msaldanha/realChain/transaction"
	"time"
	"github.com/msaldanha/realChain/keyvaluestore"
	"fmt"
	"strings"
	"github.com/msaldanha/realChain/address"
)

func assertCommonVal(val transaction.Validator, tx *transaction.Transaction) {
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

func createNonEmptyMemoryStore() *keyvaluestore.MemoryKeyValueStore {
	ms := keyvaluestore.NewMemoryKeyValueStore()
	tx := &transaction.Transaction{}
	ms.Put("genesis", tx.ToBytes())
	return ms
}

func dumpTxChain(txChain []*transaction.Transaction) {
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

func createTestAddress() *address.Address {
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