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

		block = &Block{Type: SEND, Destination: "ddddddddddddd", Previous: "ppppppppp", Signature: "ssssssss", Balance: 1,
			Work: "wwwwwwww", Account: "aaaaaaaaaa", Representative: "rrrrrrrrrrrrrrr", Source: "sssssssssssssss"}

		ok, err = bs.Store(block)
		Expect(ok).To(BeNil())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Destination not found"))

		dest := &Block{Type: OPEN, Destination: "ddddddddddddd", Previous: "ppppppppp", Signature: "ssssssss", Balance: 1,
			Work: "wwwwwwww", Account: "aaaaaaaaaa", Representative: "rrrrrrrrrrrrrrr", Source: "sssssssssssssss"}
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

		block := &Block{Type: SEND, Destination: "ddddddddddddd", Previous: "ppppppppp",
			Signature: "ce117199031c7bc2eac6658e47c4a8c3d50502375470d532b8ae451377185bf5", Balance: 1,
			Work: "wwwwwwww", Account: "aaaaaaaaaa", Representative: "rrrrrrrrrrrrrrr", Source: "sssssssssssssss"}

		dest := &Block{Type: OPEN, Destination: "ddddddddddddd", Previous: "ppppppppp", Signature: "ssssssss", Balance: 1,
			Work: "wwwwwwww", Account: "aaaaaaaaaa", Representative: "rrrrrrrrrrrrrrr", Source: "sssssssssssssss"}
		ms.Put("ddddddddddddd", dest)

		ok, err := bs.Store(block)
		Expect(ok).NotTo(BeNil())
		Expect(err).To(BeNil())

		blockFromKeyStore, _, _ := ms.Get(block.Signature)
		Expect(blockFromKeyStore).To(Equal(block))
	})

})
