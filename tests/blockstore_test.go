package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/msaldanha/realChain/block"
	"github.com/msaldanha/realChain/blockstore"
	"github.com/golang/mock/gomock"
	"github.com/msaldanha/realChain/keyvaluestore"
	"github.com/msaldanha/realChain/validator"
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
			Work: "wwwwwwww", Account: "aaaaaaaaaa", Representative: "rrrrrrrrrrrrrrr"}

		ok, err = bs.Store(block)
		Expect(ok).To(BeNil())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Destination not found"))

		dest := &Block{Type: OPEN, Link: "ddddddddddddd", Previous: "ppppppppp", Signature: "ssssssss", Balance: 1,
			Work: "wwwwwwww", Account: "aaaaaaaaaa", Representative: "rrrrrrrrrrrrrrr"}
		ms.Put("ddddddddddddd", dest)

		ok, err = bs.Store(block)
		Expect(ok).To(BeNil())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Block signature does not match"))
	})

	It("Should accept propely filled block", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		val := validator.NewBlockValidatorCreator()
		bs := blockstore.New(ms, val)

		block := &Block{Type: SEND, Link: "ddddddddddddd", Previous: "ppppppppp",
			Signature: "58853e3f4f22032c976b5569160043b04a0a0fc020d99a1a82a11051cf0598eb", Balance: 1,
			Work: "wwwwwwww", Account: "aaaaaaaaaa", Representative: "rrrrrrrrrrrrrrr"}

		dest := &Block{Type: OPEN, Link: "ddddddddddddd", Previous: "ppppppppp", Signature: "ssssssss", Balance: 1,
			Work: "wwwwwwww", Account: "aaaaaaaaaa", Representative: "rrrrrrrrrrrrrrr"}
		ms.Put("ddddddddddddd", dest)

		blk, err := bs.Store(block)
		Expect(err).To(BeNil())
		Expect(blk).NotTo(BeNil())

		blockFromKeyStore, _, _ := ms.Get(block.Signature)
		Expect(blockFromKeyStore).To(Equal(block))
	})

})
