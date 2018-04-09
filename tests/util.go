package tests

import (
	. "github.com/onsi/gomega"
	. "github.com/msaldanha/realChain/block"
	"time"
	"github.com/msaldanha/realChain/keyvaluestore"
	"fmt"
	"strings"
	"github.com/msaldanha/realChain/keypair"
	"github.com/msaldanha/realChain/address"
	"github.com/msaldanha/realChain/ledge"
)

func assertCommonVal(val Validator, block *Block) {
	ok, err := val.IsFilled(block)
	Expect(ok).To(BeFalse())
	Expect(err).NotTo(BeNil())
	Expect(err).To(Equal(ErrInvalidBlockTimestamp))

	block.Timestamp = time.Now().Unix()

	ok, err = val.IsFilled(block)
	Expect(ok).To(BeFalse())
	Expect(err).NotTo(BeNil())
	Expect(err).To(Equal(ErrBlockAccountCantBeEmpty))

	block.Account = []byte("xxxxxxxxxxxxxxxxxxx")

	ok, err = val.IsFilled(block)
	Expect(ok).To(BeFalse())
	Expect(err).NotTo(BeNil())
	Expect(err).To(Equal(ErrPreviousBlockCantBeEmpty))

	block.Previous = []byte("yyyyyyyyyyyyyyyyyyyy")

	ok, err = val.IsFilled(block)
	Expect(ok).To(BeFalse())
	Expect(err).NotTo(BeNil())
	Expect(err).To(Equal(ErrBlockSignatureCantBeEmpty))

	block.Signature = []byte("ssssssssssssssssssssss")

	ok, err = val.IsFilled(block)
	Expect(ok).To(BeFalse())
	Expect(err).NotTo(BeNil())
	Expect(err).To(Equal(ErrBlockPowNonceCantBeZero))

	block.PowNonce = 1

	ok, err = val.IsFilled(block)
	Expect(ok).To(BeFalse())
	Expect(err).NotTo(BeNil())
	Expect(err).To(Equal(ErrBlockHashCantBeEmpty))

	block.SetHash()
}

func createNonEmptyMemoryStore() *keyvaluestore.MemoryKeyValueStore {
	ms := keyvaluestore.NewMemoryKeyValueStore()
	ms.Put("genesis", &Block{})
	return ms
}

func dumpBlockChain(blockChain []*Block) {
	fmt.Println("============= Block Chain Dump start =================")
	level := 0
	for _, v := range blockChain {
		if len(v.Previous) == 0 {
			level = level + 1
		}
		fmt.Printf("%s %s %s %s %f\n", strings.Repeat("  ", level), v.Type, string(v.Account), string(v.Hash), v.Balance)
	}
	fmt.Println("============= End =================")
}

func createTestAccount() *ledge.Account {
	keys, _ := keypair.New()
	acc := &ledge.Account{Keys: keys}
	addr := address.New()
	ad, _ := addr.GenerateForKey(acc.Keys.PublicKey)
	acc.Address = string(ad)
	return acc
}