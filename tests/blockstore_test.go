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

		block = &Block{Type: SEND, Link: []byte("ddddddddddddd"), Previous: []byte("ppppppppp"), Signature: []byte("ssssssss"), Balance: 1,
			Work: []byte("wwwwwwww"), Account: []byte("aaaaaaaaaa"), Representative: []byte("rrrrrrrrrrrrrrr"), Timestamp: time.Now().Unix()}

			ok, err = bs.Store(block)
		Expect(ok).To(BeNil())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Block hash can not be empty"))

		block.SetHash()

		ok, err = bs.Store(block)
		Expect(ok).To(BeNil())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Destination not found"))

		dest := &Block{Type: OPEN, Link: []byte("ddddddddddddd"), Previous: []byte("ppppppppp"), Signature: []byte("ssssssss"), Balance: 1,
			Work: []byte("wwwwwwww"), Account: []byte("aaaaaaaaaa"), Representative: []byte("rrrrrrrrrrrrrrr"), Timestamp: time.Now().Unix()}
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

		block := &Block{Type: SEND, Link: []byte("ddddddddddddd"), Previous: []byte("ppppppppp"),
			Signature: []byte("777d713768de05cb16cbc24eef83b43b20a3a80dce05549f130aaf5a4234e4c2"), Balance: 1,
			Work: []byte("wwwwwwww"), Account: []byte("aaaaaaaaaa"), Representative: []byte("rrrrrrrrrrrrrrr"), Timestamp: 1}
		block.SetHash()

		dest := &Block{Type: OPEN, Link: []byte("ddddddddddddd"), Previous: []byte("ppppppppp"), Signature: []byte("ssssssss"), Balance: 1,
			Work: []byte("wwwwwwww"), Account: []byte("aaaaaaaaaa"), Representative: []byte("rrrrrrrrrrrrrrr"), Timestamp: 1}
		ms.Put("ddddddddddddd", dest)

		blk, err := bs.Store(block)
		Expect(err).To(BeNil())
		Expect(blk).NotTo(BeNil())

		blockFromKeyStore, _, _ := ms.Get(string(block.Hash))
		Expect(blockFromKeyStore).To(Equal(block))
	})

	It("Should calculate the PoW for the block", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		val := validator.NewBlockValidatorCreator()
		bs := blockstore.New(ms, val)

		block := &Block{Type: SEND, Link: []byte("ddddddddddddd"), Previous: []byte("ppppppppp"),
			Signature: []byte("777d713768de05cb16cbc24eef83b43b20a3a80dce05549f130aaf5a4234e4c2"), Balance: 1,
			Work: []byte("wwwwwwww"), Account: []byte("aaaaaaaaaa"), Representative: []byte("rrrrrrrrrrrrrrr"), Timestamp: 1}
		block.SetHash()

		nonce, pow, err := bs.CalculatePow(block)
		Expect(err).To(BeNil())

		Expect(nonce).To(Equal(int64(44416)))
		Expect(pow).To(Equal([]byte("000082d10e5d62e38e4fd5bb109d8809e63d546cec6900f18bb0b08d842fc9a1")))
	})

	It("Should verify the PoW for the block", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		val := validator.NewBlockValidatorCreator()
		bs := blockstore.New(ms, val)

		block := &Block{Type: SEND, Link: []byte("ddddddddddddd"), Previous: []byte("ppppppppp"),
			Signature: []byte("777d713768de05cb16cbc24eef83b43b20a3a80dce05549f130aaf5a4234e4c2"), Balance: 1,
			Work: []byte("wwwwwwww"), Account: []byte("aaaaaaaaaa"), Representative: []byte("rrrrrrrrrrrrrrr"), Timestamp: 1}
		block.SetHash()

		block.PowNonce = int64(44416)
		block.Hash = []byte("000082d10e5d62e38e4fd5bb109d8809e63d546cec6900f18bb0b08d842fc9a1")

		ok, err := bs.VerifyPow(block)
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())

	})

})
