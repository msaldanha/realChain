package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/msaldanha/realChain/block"
	"github.com/msaldanha/realChain/blockstore"
	"github.com/golang/mock/gomock"
	"github.com/msaldanha/realChain/keyvaluestore"
	"github.com/msaldanha/realChain/validator"
	"time"
)

var _ = Describe("BlockStore", func() {

	It("Should not accept empty/partially filled block", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		val := validator.NewBlockValidatorCreator()
		bs := blockstore.New(ms, val)
		block := &Block{}
		ok, err := bs.Store(block)
		Expect(ok).To(BeNil())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Invalid block type"))

		block = &Block{Type: SEND, Link: "ddddddddddddd", Previous: "ppppppppp", Signature: "ssssssss", Balance: 1,
			Work: "wwwwwwww", Account: "aaaaaaaaaa", Representative: "rrrrrrrrrrrrrrr", Timestamp: time.Now().Unix()}

		ok, err = bs.Store(block)
		Expect(ok).To(BeNil())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Destination not found"))

		dest := &Block{Type: OPEN, Link: "ddddddddddddd", Previous: "ppppppppp", Signature: "ssssssss", Balance: 1,
			Work: "wwwwwwww", Account: "aaaaaaaaaa", Representative: "rrrrrrrrrrrrrrr", Timestamp: time.Now().Unix()}
		ms.Put("ddddddddddddd", dest)

		ok, err = bs.Store(block)
		Expect(ok).To(BeNil())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Block signature does not match"))
	})

	It("Should accept properly filled block", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		val := validator.NewBlockValidatorCreator()
		bs := blockstore.New(ms, val)

		block := &Block{Type: SEND, Link: "ddddddddddddd", Previous: "ppppppppp",
			Signature: "8d2e875bd67e70a5930c60bd7e8fed1b364aa74bd5a57df7e0bad49d55558ba9", Balance: 1,
			Work: "wwwwwwww", Account: "aaaaaaaaaa", Representative: "rrrrrrrrrrrrrrr", Timestamp: 1}

		dest := &Block{Type: OPEN, Link: "ddddddddddddd", Previous: "ppppppppp", Signature: "ssssssss", Balance: 1,
			Work: "wwwwwwww", Account: "aaaaaaaaaa", Representative: "rrrrrrrrrrrrrrr", Timestamp: 1}
		ms.Put("ddddddddddddd", dest)

		blk, err := bs.Store(block)
		Expect(err).To(BeNil())
		Expect(blk).NotTo(BeNil())

		blockFromKeyStore, _, _ := ms.Get(block.Signature)
		Expect(blockFromKeyStore).To(Equal(block))
	})

})
