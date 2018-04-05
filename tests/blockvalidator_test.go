 package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/golang/mock/gomock"
	. "github.com/msaldanha/realChain/block"
	"github.com/msaldanha/realChain/keyvaluestore"
	"github.com/msaldanha/realChain/validator"
)


func assertCommonVal(val validator.BlockValidator, block *Block) {
	ok, err := val.IsFilled(block)
	Expect(ok).To(BeFalse())
	Expect(err).NotTo(BeNil())
	Expect(err.Error()).To(Equal("Previous block can not be empty"))

	block.Previous = "yyyyyyyyyyyyyyyyyyyy"

	ok, err = val.IsFilled(block)
	Expect(ok).To(BeFalse())
	Expect(err).NotTo(BeNil())
	Expect(err.Error()).To(Equal("Block signature can not be empty"))

	block.Signature = "ssssssssssssssssssssss"

	ok, err = val.IsFilled(block)
	Expect(ok).To(BeFalse())
	Expect(err).NotTo(BeNil())
	Expect(err.Error()).To(Equal("Block PoW can not be empty"))

	block.Work = "ssssssssssssssssssssss"

	ok, err = val.IsFilled(block)
	Expect(ok).To(BeTrue())
	Expect(err).To(BeNil())
}

var _ = Describe("BlockValidator", func() {
	It("Should not accept empty/partial filled block for OPEN type", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		val := validator.NewBlockValidatorCreator().CreateValidatorForBlock(OPEN, ms)

		block := &Block{}
		ok, err := val.IsFilled(block)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Invalid block type"))

		block.Type = OPEN

		ok, err = val.IsFilled(block)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Block account can not be empty"))

		block.Account = "xxxxxxxxxxxxxxxxxxx"

		ok, err = val.IsFilled(block)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Block representative can not be empty"))

		block.Representative = "xxxxxxxxxxxxxxxxxxx"

		ok, err = val.IsFilled(block)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Block signature can not be empty"))

		block.Signature = "ssssssssssssssssssssss"

		ok, err = val.IsFilled(block)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Block PoW can not be empty"))

		block.Work = "ssssssssssssssssssssss"

		ok, err = val.IsFilled(block)
		Expect(ok).To(BeTrue())
		Expect(err).To(BeNil())
	})

	It("Should not accept empty/partial filled block for SEND type", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		val := validator.NewBlockValidatorCreator().CreateValidatorForBlock(SEND, ms)

		block := &Block{}
		ok, err := val.IsFilled(block)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Invalid block type"))

		block.Type = SEND

		ok, err = val.IsFilled(block)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Block destination can not be empty"))

		block.Link = "xxxxxxxxxxxxxxxxxxx"

		assertCommonVal(val, block)
	})

	It("Should not accept empty/partial filled block for RECEIVE type", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		val := validator.NewBlockValidatorCreator().CreateValidatorForBlock(RECEIVE, ms)

		block := &Block{}
		ok, err := val.IsFilled(block)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Invalid block type"))

		block.Type = RECEIVE

		ok, err = val.IsFilled(block)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Block source can not be empty"))

		block.Link = "xxxxxxxxxxxxxxxxxxx"

		assertCommonVal(val, block)
	})

	It("Should not accept empty/partial filled block for CHANGE type", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		val := validator.NewBlockValidatorCreator().CreateValidatorForBlock(CHANGE, ms)

		block := &Block{}
		ok, err := val.IsFilled(block)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Invalid block type"))

		block.Type = CHANGE

		ok, err = val.IsFilled(block)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Block representative can not be empty"))

		block.Representative = "xxxxxxxxxxxxxxxxxxx"

		assertCommonVal(val, block)
	})

	It("Should not accept OPEN block with invalid fields", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		val := validator.NewBlockValidatorCreator().CreateValidatorForBlock(OPEN, ms)
		block := &Block{Type: OPEN, Link: "ddddddddddddd", Previous: "ppppppppp",
		Signature: "ssssssss", Balance: 1,
			Work: "wwwwwwww", Account:"aaaaaaaaaa", Representative:"rrrrrrrrrrrrrrr"}

		ok, err := val.IsValid(block)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Block signature does not match"))

		block.Signature = "d8ab2722a816681533a434728f0bd9810e9e0827a4357d25b6c1d8a6efb12379"

		ok, err = val.IsValid(block)
		Expect(ok).To(BeTrue())
		Expect(err).To(BeNil())
	})

	It("Should not accept SEND block with invalid fields", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		val := validator.NewBlockValidatorCreator().CreateValidatorForBlock(SEND, ms)
		block := &Block{Type: SEND, Link: "ddddddddddddd", Previous: "ppppppppp", Signature: "ssssssss", Balance: 1,
			Work: "wwwwwwww", Account:"aaaaaaaaaa", Representative:"rrrrrrrrrrrrrrr"}

		ok, err := val.IsValid(block)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Destination not found"))

		dest := &Block{Type: OPEN, Link: "ddddddddddddd", Previous: "ppppppppp", Signature: "ssssssss", Balance: 1,
			Work: "wwwwwwww", Account:"aaaaaaaaaa", Representative:"rrrrrrrrrrrrrrr"}
		ms.Put("ddddddddddddd", dest)

		ok, err = val.IsValid(block)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Block signature does not match"))

		block.Signature = "58853e3f4f22032c976b5569160043b04a0a0fc020d99a1a82a11051cf0598eb"

		ok, err = val.IsValid(block)
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())

	})

	It("Should not accept RECEIVE block with invalid fields", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		val := validator.NewBlockValidatorCreator().CreateValidatorForBlock(RECEIVE, ms)
		block := &Block{Type: RECEIVE, Link: "ddddddddddddd", Previous: "ppppppppp", Signature: "ssssssss", Balance: 1,
			Work: "wwwwwwww", Account:"aaaaaaaaaa", Representative:"rrrrrrrrrrrrrrr"}

		ok, err := val.IsValid(block)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Source not found"))

		source := &Block{Type: OPEN, Link: "ddddddddddddd", Previous: "ppppppppp", Signature: "ssssssss", Balance: 1,
			Work: "wwwwwwww", Account:"aaaaaaaaaa", Representative:"rrrrrrrrrrrrrrr"}
		ms.Put("ddddddddddddd", source)

		ok, err = val.IsValid(block)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Source of invalid type"))

		source.Type = SEND

		ok, err = val.IsValid(block)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Block signature does not match"))

		block.Signature = "33539bf26a3be64bf19ba4f62c850c3dd0a1d70072c9fccd3b1b8e9441cb00c9"

		ok, err = val.IsValid(block)
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())
	})

	It("Should not accept CHANGE block with invalid fields", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		val := validator.NewBlockValidatorCreator().CreateValidatorForBlock(CHANGE, ms)
		block := &Block{Type: CHANGE, Link: "ddddddddddddd", Previous: "ppppppppp", Signature: "ssssssss", Balance: 1,
			Work: "wwwwwwww", Account:"aaaaaaaaaa", Representative:"rrrrrrrrrrrrrrr"}

		ok, err := val.IsValid(block)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Block signature does not match"))

		block.Signature = "76f4d115e9ab6870f70aa54850dbd46d01f3901edd384944f4a7ad1db7ec6550"

		ok, err = val.IsValid(block)
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())
	})

})
