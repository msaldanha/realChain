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

	block.Previous = []byte("yyyyyyyyyyyyyyyyyyyy")

	ok, err = val.IsFilled(block)
	Expect(ok).To(BeFalse())
	Expect(err).NotTo(BeNil())
	Expect(err.Error()).To(Equal("Block signature can not be empty"))

	block.Signature = []byte("ssssssssssssssssssssss")

	ok, err = val.IsFilled(block)
	Expect(ok).To(BeFalse())
	Expect(err).NotTo(BeNil())
	Expect(err.Error()).To(Equal("Block PoW can not be empty"))

	block.Work = []byte("ssssssssssssssssssssss")

	ok, err = val.IsFilled(block)
	Expect(ok).To(BeFalse())
	Expect(err).NotTo(BeNil())
	Expect(err.Error()).To(Equal("Block hash can not be empty"))

	block.SetHash()

	ok, err = val.IsFilled(block)
	Expect(err).To(BeNil())
	Expect(ok).To(BeTrue())
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

		block.Account = []byte("xxxxxxxxxxxxxxxxxxx")

		ok, err = val.IsFilled(block)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Block representative can not be empty"))

		block.Representative = []byte("xxxxxxxxxxxxxxxxxxx")

		ok, err = val.IsFilled(block)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Block signature can not be empty"))

		block.Signature = []byte("ssssssssssssssssssssss")

		ok, err = val.IsFilled(block)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Block PoW can not be empty"))

		block.Work = []byte("ssssssssssssssssssssss")

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

		block.Link = []byte("xxxxxxxxxxxxxxxxxxx")

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

		block.Link = []byte("xxxxxxxxxxxxxxxxxxx")

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

		block.Representative = []byte("xxxxxxxxxxxxxxxxxxx")

		assertCommonVal(val, block)
	})

	It("Should not accept OPEN block with invalid fields", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		val := validator.NewBlockValidatorCreator().CreateValidatorForBlock(OPEN, ms)
		block := &Block{Type: OPEN, Link: []byte([]byte("ddddddddddddd")), Previous: []byte([]byte("ppppppppp")),
		Signature: []byte([]byte("ssssssss")), Balance: 1, Timestamp: 1,
			Work: []byte("wwwwwwww"), Account:[]byte("aaaaaaaaaa"), Representative:[]byte("rrrrrrrrrrrrrrr")}
		block.SetHash()

		ok, err := val.IsValid(block)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Block signature does not match"))

		block.Signature = []byte("7188ec0ca01d31f6e181b4cfd21c7830f66b0c4966b7ecd09431c62f45f4e504")

		ok, err = val.IsValid(block)
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())
	})

	It("Should not accept SEND block with invalid fields", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		val := validator.NewBlockValidatorCreator().CreateValidatorForBlock(SEND, ms)
		block := &Block{Type: SEND, Link: []byte("ddddddddddddd"), Previous: []byte("ppppppppp"), Signature: []byte("ssssssss"), Balance: 1,
			Work: []byte("wwwwwwww"), Account:[]byte("aaaaaaaaaa"), Representative:[]byte("rrrrrrrrrrrrrrr"), Timestamp: 1}
		block.SetHash()

		ok, err := val.IsValid(block)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Destination not found"))

		dest := &Block{Type: OPEN, Link: []byte("ddddddddddddd"), Previous: []byte("ppppppppp"), Signature: []byte("ssssssss"), Balance: 1,
			Work: []byte("wwwwwwww"), Account:[]byte("aaaaaaaaaa"), Representative:[]byte("rrrrrrrrrrrrrrr"), Timestamp: 1}
		ms.Put("ddddddddddddd", dest)

		ok, err = val.IsValid(block)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Block signature does not match"))

		block.Signature = []byte("777d713768de05cb16cbc24eef83b43b20a3a80dce05549f130aaf5a4234e4c2")

		ok, err = val.IsValid(block)
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())

	})

	It("Should not accept RECEIVE block with invalid fields", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		val := validator.NewBlockValidatorCreator().CreateValidatorForBlock(RECEIVE, ms)
		block := &Block{Type: RECEIVE, Link: []byte("ddddddddddddd"), Previous: []byte("ppppppppp"), Signature: []byte("ssssssss"), Balance: 1,
			Work: []byte("wwwwwwww"), Account:[]byte("aaaaaaaaaa"), Representative:[]byte("rrrrrrrrrrrrrrr"), Timestamp: 1}
		block.SetHash()

		ok, err := val.IsValid(block)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Source not found"))

		source := &Block{Type: OPEN, Link: []byte("ddddddddddddd"), Previous: []byte("ppppppppp"), Signature: []byte("ssssssss"), Balance: 1,
			Work: []byte("wwwwwwww"), Account:[]byte("aaaaaaaaaa"), Representative:[]byte("rrrrrrrrrrrrrrr"), Timestamp: 1}
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

		block.Signature = []byte("048dd436a5cacb7c0bcf078ce7234de77a0f23b09f531831820f3ff081a62421")

		ok, err = val.IsValid(block)
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())
	})

	It("Should not accept CHANGE block with invalid fields", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		val := validator.NewBlockValidatorCreator().CreateValidatorForBlock(CHANGE, ms)
		block := &Block{Type: CHANGE, Link: []byte("ddddddddddddd"), Previous: []byte("ppppppppp"), Signature: []byte("ssssssss"), Balance: 1,
			Work: []byte("wwwwwwww"), Account:[]byte("aaaaaaaaaa"), Representative:[]byte("rrrrrrrrrrrrrrr"), Timestamp: 1}
		block.SetHash()

		ok, err := val.IsValid(block)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Block signature does not match"))

		block.Signature = []byte("e39f97b4322b4016412fbb2d99358c01d01429afeca8338af116d28acaf40d99")

		ok, err = val.IsValid(block)
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())
	})

})
