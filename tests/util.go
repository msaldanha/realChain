package tests

import (
	. "github.com/onsi/gomega"
	. "github.com/msaldanha/realChain/transaction"
	"time"
	"github.com/msaldanha/realChain/keyvaluestore"
	"fmt"
	"strings"
	"github.com/msaldanha/realChain/keypair"
	"github.com/msaldanha/realChain/address"
	"github.com/msaldanha/realChain/ledger"
)

func assertCommonVal(val Validator, tx *Transaction) {
	ok, err := val.IsFilled(tx)
	Expect(ok).To(BeFalse())
	Expect(err).NotTo(BeNil())
	Expect(err).To(Equal(ErrInvalidTransactionTimestamp))

	tx.Timestamp = time.Now().Unix()

	ok, err = val.IsFilled(tx)
	Expect(ok).To(BeFalse())
	Expect(err).NotTo(BeNil())
	Expect(err).To(Equal(ErrTransactionAccountCantBeEmpty))

	tx.Account = []byte("xxxxxxxxxxxxxxxxxxx")

	ok, err = val.IsFilled(tx)
	Expect(ok).To(BeFalse())
	Expect(err).NotTo(BeNil())
	Expect(err).To(Equal(ErrPreviousTransactionCantBeEmpty))

	tx.Previous = []byte("yyyyyyyyyyyyyyyyyyyy")

	ok, err = val.IsFilled(tx)
	Expect(ok).To(BeFalse())
	Expect(err).NotTo(BeNil())
	Expect(err).To(Equal(ErrTransactionSignatureCantBeEmpty))

	tx.Signature = []byte("ssssssssssssssssssssss")

	ok, err = val.IsFilled(tx)
	Expect(ok).To(BeFalse())
	Expect(err).NotTo(BeNil())
	Expect(err).To(Equal(ErrTransactionPowNonceCantBeZero))

	tx.PowNonce = 1

	ok, err = val.IsFilled(tx)
	Expect(ok).To(BeFalse())
	Expect(err).NotTo(BeNil())
	Expect(err).To(Equal(ErrTransactionHashCantBeEmpty))

	tx.SetHash()

	ok, err = val.IsFilled(tx)
	Expect(ok).To(BeFalse())
	Expect(err).NotTo(BeNil())
	Expect(err).To(Equal(ErrPubKeyCantBeEmpty))

	tx.PubKey = []byte("ssssssssssssssssssssss")
}

func createNonEmptyMemoryStore() *keyvaluestore.MemoryKeyValueStore {
	ms := keyvaluestore.NewMemoryKeyValueStore()
	tx := &Transaction{}
	ms.Put("genesis", tx.ToBytes())
	return ms
}

func dumpTxChain(txChain []*Transaction) {
	fmt.Println("============= Transaction Chain Dump start =================")
	level := 0
	for _, v := range txChain {
		if len(v.Previous) == 0 {
			level = level + 1
		}
		fmt.Printf("%s %s %s %s %s %f\n", strings.Repeat("  ", level), v.Type, string(v.Account), string(v.Hash), string(v.Previous), v.Balance)
	}
	fmt.Println("============= End =================")
}

func createTestAccount() *ledger.Account {
	keys, _ := keypair.New()
	acc := &ledger.Account{Keys: keys}
	addr := address.New()
	ad, _ := addr.GenerateForKey(acc.Keys.PublicKey)
	acc.Address = string(ad)
	return acc
}