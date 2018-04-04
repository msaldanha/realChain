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
		Expect(err.Error()).To(Equal("Block source can not be empty"))

		block.Source = "xxxxxxxxxxxxxxxxxxx"

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

		block.Destination = "xxxxxxxxxxxxxxxxxxx"

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

		block.Source = "xxxxxxxxxxxxxxxxxxx"

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
		block := &Block{Type: OPEN, Destination: "ddddddddddddd", Previous: "ppppppppp", Signature: "ssssssss", Balance: 1,
			Work: "wwwwwwww", Account:"aaaaaaaaaa", Representative:"rrrrrrrrrrrrrrr", Source:"sssssssssssssss"}

		ok, err := val.IsValid(block)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Source not found"))

		source := &Block{Type: OPEN, Destination: "ddddddddddddd", Previous: "ppppppppp", Signature: "ssssssss", Balance: 1,
			Work: "wwwwwwww", Account:"aaaaaaaaaa", Representative:"rrrrrrrrrrrrrrr", Source:"sssssssssssssss"}
		ms.Put("sssssssssssssss", source)

		ok, err = val.IsValid(block)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Source of invalid type"))

		source.Type = SEND

		ok, err = val.IsValid(block)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Block signature does not match"))

		block.Signature = "eaba3aabc0e9f1a0f3a2590e6fad7375609ad7596a6dbed0ebe71c8cf6b8004d"

		ok, err = val.IsValid(block)
		Expect(ok).To(BeTrue())
		Expect(err).To(BeNil())
	})

	It("Should not accept SEND block with invalid fields", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		val := validator.NewBlockValidatorCreator().CreateValidatorForBlock(SEND, ms)
		block := &Block{Type: SEND, Destination: "ddddddddddddd", Previous: "ppppppppp", Signature: "ssssssss", Balance: 1,
			Work: "wwwwwwww", Account:"aaaaaaaaaa", Representative:"rrrrrrrrrrrrrrr", Source:"sssssssssssssss"}

		ok, err := val.IsValid(block)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Destination not found"))

		dest := &Block{Type: OPEN, Destination: "ddddddddddddd", Previous: "ppppppppp", Signature: "ssssssss", Balance: 1,
			Work: "wwwwwwww", Account:"aaaaaaaaaa", Representative:"rrrrrrrrrrrrrrr", Source:"sssssssssssssss"}
		ms.Put("ddddddddddddd", dest)

		ok, err = val.IsValid(block)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Block signature does not match"))

		block.Signature = "ce117199031c7bc2eac6658e47c4a8c3d50502375470d532b8ae451377185bf5"

		ok, err = val.IsValid(block)
		Expect(ok).To(BeTrue())
		Expect(err).To(BeNil())
	})

	It("Should not accept RECEIVE block with invalid fields", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		val := validator.NewBlockValidatorCreator().CreateValidatorForBlock(RECEIVE, ms)
		block := &Block{Type: RECEIVE, Destination: "ddddddddddddd", Previous: "ppppppppp", Signature: "ssssssss", Balance: 1,
			Work: "wwwwwwww", Account:"aaaaaaaaaa", Representative:"rrrrrrrrrrrrrrr", Source:"sssssssssssssss"}

		ok, err := val.IsValid(block)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Source not found"))

		source := &Block{Type: OPEN, Destination: "ddddddddddddd", Previous: "ppppppppp", Signature: "ssssssss", Balance: 1,
			Work: "wwwwwwww", Account:"aaaaaaaaaa", Representative:"rrrrrrrrrrrrrrr", Source:"sssssssssssssss"}
		ms.Put("sssssssssssssss", source)

		ok, err = val.IsValid(block)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Source of invalid type"))

		source.Type = SEND

		ok, err = val.IsValid(block)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Block signature does not match"))

		block.Signature = "57eff6430247760173ca8072cd18fd01a0629fa8ac164a624dfe89044f596ff6"

		ok, err = val.IsValid(block)
		Expect(ok).To(BeTrue())
		Expect(err).To(BeNil())
	})

	It("Should not accept CHANGE block with invalid fields", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		val := validator.NewBlockValidatorCreator().CreateValidatorForBlock(CHANGE, ms)
		block := &Block{Type: CHANGE, Destination: "ddddddddddddd", Previous: "ppppppppp", Signature: "ssssssss", Balance: 1,
			Work: "wwwwwwww", Account:"aaaaaaaaaa", Representative:"rrrrrrrrrrrrrrr", Source:"sssssssssssssss"}

		ok, err := val.IsValid(block)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Block signature does not match"))

		block.Signature = "e280768734004edfb354fa762dd3a311e5c3f1f23d506c369bbd4ae35f73686c"

		ok, err = val.IsValid(block)
		Expect(ok).To(BeTrue())
		Expect(err).To(BeNil())
	})

})
