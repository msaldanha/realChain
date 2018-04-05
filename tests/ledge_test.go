package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/msaldanha/realChain/blockstore"
	"github.com/golang/mock/gomock"
	"github.com/msaldanha/realChain/keyvaluestore"
	"github.com/msaldanha/realChain/validator"
	"github.com/msaldanha/realChain/ledge"
)

var _ = Describe("Ledge", func() {
	It("Should create the Genesis block", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		val := validator.NewBlockValidatorCreator()
		bs := blockstore.New(ms, val)

		ld := ledge.New()
		ld.Use(bs)

		blk, err := ld.Initialize(1000)
		Expect(err).To(BeNil())

		blk, err = bs.Retrieve(string(blk.Hash))
		Expect(err).To(BeNil())
		Expect(blk).NotTo(BeNil())
	})

	It("Should NOT create the Genesis block twice", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		val := validator.NewBlockValidatorCreator()
		bs := blockstore.New(ms, val)

		ld := ledge.New()
		ld.Use(bs)

		blk, err := ld.Initialize(1000)
		Expect(err).To(BeNil())
		Expect(blk).NotTo(BeNil())

		blk, err = ld.Initialize(1000)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Ledge already initialized"))
	})
})