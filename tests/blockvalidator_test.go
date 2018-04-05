 package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/golang/mock/gomock"
	. "github.com/msaldanha/realChain/block"
	"github.com/msaldanha/realChain/keyvaluestore"
	"github.com/msaldanha/realChain/validator"
	"time"
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
		Signature: "ssssssss", Balance: 1, Timestamp: 1,
			Work: "wwwwwwww", Account:"aaaaaaaaaa", Representative:"rrrrrrrrrrrrrrr"}

		ok, err := val.IsValid(block)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Block signature does not match"))

		block.Signature = "e0f7b968d899426dd4372b507374e27de643d20588e0a40875d1a0cc4350e431"

		ok, err = val.IsValid(block)
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())
	})

	It("Should not accept SEND block with invalid fields", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		val := validator.NewBlockValidatorCreator().CreateValidatorForBlock(SEND, ms)
		block := &Block{Type: SEND, Link: "ddddddddddddd", Previous: "ppppppppp", Signature: "ssssssss", Balance: 1,
			Work: "wwwwwwww", Account:"aaaaaaaaaa", Representative:"rrrrrrrrrrrrrrr", Timestamp: 1}

		ok, err := val.IsValid(block)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Destination not found"))

		dest := &Block{Type: OPEN, Link: "ddddddddddddd", Previous: "ppppppppp", Signature: "ssssssss", Balance: 1,
			Work: "wwwwwwww", Account:"aaaaaaaaaa", Representative:"rrrrrrrrrrrrrrr", Timestamp: 1}
		ms.Put("ddddddddddddd", dest)

		ok, err = val.IsValid(block)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Block signature does not match"))

		block.Signature = "8d2e875bd67e70a5930c60bd7e8fed1b364aa74bd5a57df7e0bad49d55558ba9"

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
			Work: "wwwwwwww", Account:"aaaaaaaaaa", Representative:"rrrrrrrrrrrrrrr", Timestamp: 1}

		ok, err := val.IsValid(block)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Source not found"))

		source := &Block{Type: OPEN, Link: "ddddddddddddd", Previous: "ppppppppp", Signature: "ssssssss", Balance: 1,
			Work: "wwwwwwww", Account:"aaaaaaaaaa", Representative:"rrrrrrrrrrrrrrr", Timestamp: 1}
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

		block.Signature = "df8c4b90098e7d69aa3623a5b8aa646d25f23634b1013a8cd5b131d9b86dfce0"

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
			Work: "wwwwwwww", Account:"aaaaaaaaaa", Representative:"rrrrrrrrrrrrrrr", Timestamp: 1}

		ok, err := val.IsValid(block)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Block signature does not match"))

		block.Signature = "e38dc06416c781dd6c627c03a44ad954073a15a91d6e51a56f292a02e48fd449"

		ok, err = val.IsValid(block)
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())
	})

})
