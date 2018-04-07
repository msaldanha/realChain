package tests

import (
	. "github.com/onsi/gomega"
	. "github.com/msaldanha/realChain/block"
	"github.com/msaldanha/realChain/validator"
	"time"
	"github.com/msaldanha/realChain/keyvaluestore"
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