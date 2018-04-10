package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/golang/mock/gomock"
	"github.com/msaldanha/realChain/block"
	"time"
)

var _ = Describe("Validator", func() {
	It("Should not accept empty/partial filled block for OPEN type", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := createNonEmptyMemoryStore()
		val := block.NewBlockValidatorCreator().CreateValidatorForBlock(block.OPEN, ms)

		blk := &block.Block{}
		ok, err := val.IsFilled(blk)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(block.ErrInvalidBlockType))

		blk.Type = block.OPEN

		ok, err = val.IsFilled(blk)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(block.ErrInvalidBlockTimestamp))

		blk.Timestamp = time.Now().Unix()

		ok, err = val.IsFilled(blk)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(block.ErrBlockAccountCantBeEmpty))

		blk.Account = []byte("xxxxxxxxxxxxxxxxxxx")

		ok, err = val.IsFilled(blk)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(block.ErrBlockSignatureCantBeEmpty))

		blk.Signature = []byte("ssssssssssssssssssssss")

		ok, err = val.IsFilled(blk)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(block.ErrBlockPowNonceCantBeZero))

		blk.PowNonce = 1

		ok, err = val.IsFilled(blk)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(block.ErrBlockHashCantBeEmpty))

		blk.Hash = []byte("xxxxxxxxxxxxxxxxxxx")

		ok, err = val.IsFilled(blk)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(block.ErrPubKeyCantBeEmpty))

		blk.PubKey = []byte("xxxxxxxxxxxxxxxxxxx")

		ok, err = val.IsFilled(blk)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(block.ErrBlockLinkCantBeEmpty))

		blk.Link = []byte("xxxxxxxxxxxxxxxxxxx")

		ok, err = val.IsFilled(blk)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(block.ErrBlockRepresentativeCantBeEmpty))

		blk.Representative = []byte("xxxxxxxxxxxxxxxxxxx")

		ok, err = val.IsFilled(blk)
		Expect(ok).To(BeTrue())
		Expect(err).To(BeNil())
	})

	It("Should not accept empty/partial filled block for SEND type", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := createNonEmptyMemoryStore()
		val := block.NewBlockValidatorCreator().CreateValidatorForBlock(block.SEND, ms)

		blk := &block.Block{}
		ok, err := val.IsFilled(blk)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(block.ErrInvalidBlockType))

		blk.Type = block.SEND

		assertCommonVal(val, blk)

		ok, err = val.IsFilled(blk)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(block.ErrBlockLinkCantBeEmpty))

		blk.Link = []byte("xxxxxxxxxxxxxxxxxxx")
	})

	It("Should not accept empty/partial filled block for RECEIVE type", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := createNonEmptyMemoryStore()
		val := block.NewBlockValidatorCreator().CreateValidatorForBlock(block.RECEIVE, ms)

		blk := &block.Block{}
		ok, err := val.IsFilled(blk)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(block.ErrInvalidBlockType))

		blk.Type = block.RECEIVE

		assertCommonVal(val, blk)

		ok, err = val.IsFilled(blk)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(block.ErrBlockLinkCantBeEmpty))

		blk.Link = []byte("xxxxxxxxxxxxxxxxxxx")
	})

	It("Should not accept empty/partial filled block for CHANGE type", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := createNonEmptyMemoryStore()
		val := block.NewBlockValidatorCreator().CreateValidatorForBlock(block.CHANGE, ms)

		blk := &block.Block{}
		ok, err := val.IsFilled(blk)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(block.ErrInvalidBlockType))

		blk.Type = block.CHANGE

		assertCommonVal(val, blk)

		ok, err = val.IsFilled(blk)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(block.ErrBlockRepresentativeCantBeEmpty))

		blk.Representative = []byte("xxxxxxxxxxxxxxxxxxx")
	})

	It("Should not accept OPEN block with invalid fields", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := createNonEmptyMemoryStore()
		val := block.NewBlockValidatorCreator().CreateValidatorForBlock(block.OPEN, ms)
		blk := &block.Block{Type: block.OPEN, Link: []byte([]byte("ddddddddddddd")), Previous: []byte([]byte("ppppppppp")),
		Signature: []byte([]byte("ssssssss")), Balance: 1, Timestamp: 1, PubKey: []byte("kkkkkkk"),
			PowNonce: 1, Representative:[]byte("rrrrrrrrrrrrrrr")}
		blk.SetHash()

		ok, err := val.IsValid(blk)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(block.ErrBlockAccountCantBeEmpty))

		blk.Account = []byte("aaaaaaaaaa")

		ok, err = val.IsValid(blk)
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())
	})

	It("Should not accept SEND block with invalid fields", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := createNonEmptyMemoryStore()
		val := block.NewBlockValidatorCreator().CreateValidatorForBlock(block.SEND, ms)
		blk := &block.Block{Type: block.SEND, Previous: []byte("ppppppppp"), Signature: []byte("ssssssss"), Balance: 1,
			PowNonce: 1, Account:[]byte("aaaaaaaaaa"), Representative:[]byte("rrrrrrrrrrrrrrr"), Timestamp: 1, PubKey: []byte("kkkkkkk")}
		blk.SetHash()

		ok, err := val.IsValid(blk)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(block.ErrBlockLinkCantBeEmpty))

		blk.Link = []byte("ddddddddddddd")

		ok, err = val.IsValid(blk)
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())

	})

	It("Should not accept RECEIVE block with invalid fields", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := createNonEmptyMemoryStore()
		val := block.NewBlockValidatorCreator().CreateValidatorForBlock(block.RECEIVE, ms)
		blk := &block.Block{Type: block.RECEIVE, Link: []byte("ddddddddddddd"), Previous: []byte("ppppppppp"), Signature: []byte("ssssssss"), Balance: 1,
			PowNonce: 1, Account:[]byte("aaaaaaaaaa"), Representative:[]byte("rrrrrrrrrrrrrrr"), Timestamp: 1, PubKey: []byte("kkkkkkk")}
		blk.SetHash()

		ok, err := val.IsValid(blk)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(block.ErrSourceNotFound))

		source := &block.Block{Type: block.OPEN, Link: []byte("ddddddddddddd"), Previous: []byte("ppppppppp"), Signature: []byte("ssssssss"), Balance: 1,
			PowNonce: 1, Account:[]byte("aaaaaaaaaa"), Representative:[]byte("rrrrrrrrrrrrrrr"), Timestamp: 1, PubKey: []byte("kkkkkkk")}
		ms.Put("ddddddddddddd", source)

		ok, err = val.IsValid(blk)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(block.ErrInvalidSourceType))

		source.Type = block.SEND

		ok, err = val.IsValid(blk)
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())
	})

	It("Should not accept CHANGE block with invalid fields", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := createNonEmptyMemoryStore()
		val := block.NewBlockValidatorCreator().CreateValidatorForBlock(block.CHANGE, ms)
		blk := &block.Block{Type: block.CHANGE, Link: []byte("ddddddddddddd"), Previous: []byte("ppppppppp"), Signature: []byte("ssssssss"), Balance: 1,
			PowNonce: 1, Account:[]byte("aaaaaaaaaa"), Representative:[]byte("rrrrrrrrrrrrrrr"), Timestamp: 1, PubKey: []byte("kkkkkkk")}
		blk.SetHash()

		ok, err := val.IsValid(blk)
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())
	})

})
