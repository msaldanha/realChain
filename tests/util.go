package tests

import (
	. "github.com/onsi/gomega"
	. "github.com/msaldanha/realChain/block"
	"github.com/msaldanha/realChain/validator"
	"time"
	"github.com/msaldanha/realChain/keyvaluestore"
	"fmt"
	"strings"
	"github.com/msaldanha/realChain/keypair"
	"github.com/msaldanha/realChain/address"
	"github.com/msaldanha/realChain/ledge"
)

func assertCommonVal(val validator.BlockValidator, block *Block) {
	ok, err := val.IsFilled(block)
	Expect(ok).To(BeFalse())
	Expect(err).NotTo(BeNil())
	Expect(err.Error()).To(Equal("Invalid block timestamp"))

	block.Timestamp = time.Now().Unix()

	ok, err = val.IsFilled(block)
	Expect(ok).To(BeFalse())
	Expect(err).NotTo(BeNil())
	Expect(err.Error()).To(Equal("Block account can not be empty"))

	block.Account = []byte("xxxxxxxxxxxxxxxxxxx")

	ok, err = val.IsFilled(block)
	Expect(ok).To(BeFalse())
	Expect(err).NotTo(BeNil())
	Expect(err.Error()).To(Equal("Previous block can not be empty"))

	block.Previous = []byte("yyyyyyyyyyyyyyyyyyyy")

	ok, err = val.IsFilled(block)
	Expect(ok).To(BeFalse())
	Expect(err).NotTo(BeNil())
	Expect(err.Error()).To(Equal("Block signature can not be empty"))

	block.Signature = []byte("ssssssssssssssssssssss")

	ok, err = val.IsFilled(block)
	Expect(ok).To(BeFalse())
	Expect(err).NotTo(BeNil())
	Expect(err.Error()).To(Equal("Block PoW nonce can not be zero"))

	block.PowNonce = 1

	ok, err = val.IsFilled(block)
	Expect(ok).To(BeFalse())
	Expect(err).NotTo(BeNil())
	Expect(err.Error()).To(Equal("Block hash can not be empty"))

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