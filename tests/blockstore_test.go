package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/msaldanha/realChain/block"
	"github.com/msaldanha/realChain/blockstore"
	"github.com/golang/mock/gomock"
	"github.com/msaldanha/realChain/validator"
	"time"
)

var _ = Describe("BlockStore", func() {

	It("Should not accept empty/partially filled block", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := createNonEmptyMemoryStore()
		val := validator.NewBlockValidatorCreator()
		bs := blockstore.New(ms, val)
		block := &Block{}
		ok, err := bs.Store(block)
		Expect(ok).To(BeNil())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Invalid block type"))

		block = &Block{Type: SEND, Link: []byte("ddddddddddddd"), Previous: []byte("ppppppppp"), Signature: []byte("ssssssss"), Balance: 1,
			PowNonce: 1, Account: []byte("aaaaaaaaaa"), Representative: []byte("rrrrrrrrrrrrrrr"), Timestamp: time.Now().Unix()}

			ok, err = bs.Store(block)
		Expect(ok).To(BeNil())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Block hash can not be empty"))

		block.SetHash()

		dest := &Block{Type: OPEN, Link: []byte("ddddddddddddd"), Previous: []byte("ppppppppp"), Signature: []byte("ssssssss"), Balance: 1,
			PowNonce: 1, Account: []byte("aaaaaaaaaa"), Representative: []byte("rrrrrrrrrrrrrrr"), Timestamp: time.Now().Unix()}
		ms.Put("ddddddddddddd", dest)

		ok, err = bs.Store(block)
		Expect(ok).To(BeNil())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Block signature does not match"))
	})

	It("Should accept properly filled block", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := createNonEmptyMemoryStore()
		val := validator.NewBlockValidatorCreator()
		bs := blockstore.New(ms, val)

		block := &Block{Type: SEND, Link: []byte("ddddddddddddd"), Previous: []byte("ppppppppp"),
			Signature: []byte("a246ce6b1d2b57ac33073127d8f9539fca32fb48481d46d734bf3308796ee18b"), Balance: 1,
			PowNonce: 1, Account: []byte("aaaaaaaaaa"), Representative: []byte("rrrrrrrrrrrrrrr"), Timestamp: 1}
		block.SetHash()

		dest := &Block{Type: OPEN, Link: []byte("ddddddddddddd"), Previous: []byte("ppppppppp"), Signature: []byte("ssssssss"), Balance: 1,
			PowNonce: 1, Account: []byte("aaaaaaaaaa"), Representative: []byte("rrrrrrrrrrrrrrr"), Timestamp: 1}
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

		ms := createNonEmptyMemoryStore()
		val := validator.NewBlockValidatorCreator()
		bs := blockstore.New(ms, val)

		block := &Block{Type: SEND, Link: []byte("ddddddddddddd"), Previous: []byte("ppppppppp"),
			Signature: []byte("777d713768de05cb16cbc24eef83b43b20a3a80dce05549f130aaf5a4234e4c2"), Balance: 1,
			PowNonce: 1, Account: []byte("aaaaaaaaaa"), Representative: []byte("rrrrrrrrrrrrrrr"), Timestamp: 1}
		block.SetHash()

		nonce, pow, err := bs.CalculatePow(block)
		Expect(err).To(BeNil())

		Expect(nonce).To(Equal(int64(33794)))
		Expect(pow).To(Equal([]byte("0000f4722f6416ddb43a4ee56921dd3a24c93b051a570e14ca07cd174517cf12")))
	})

	It("Should verify the PoW for the block", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := createNonEmptyMemoryStore()
		val := validator.NewBlockValidatorCreator()
		bs := blockstore.New(ms, val)

		block := &Block{Type: SEND, Link: []byte("ddddddddddddd"), Previous: []byte("ppppppppp"),
			Signature: []byte("777d713768de05cb16cbc24eef83b43b20a3a80dce05549f130aaf5a4234e4c2"), Balance: 1,
			PowNonce: 1, Account: []byte("aaaaaaaaaa"), Representative: []byte("rrrrrrrrrrrrrrr"), Timestamp: 1}
		block.SetHash()

		block.PowNonce = int64(33794)
		block.Hash = []byte("0000f4722f6416ddb43a4ee56921dd3a24c93b051a570e14ca07cd174517cf12")

		ok, err := bs.VerifyPow(block)
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())

	})

	It("Should extract the chain for the block", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := createNonEmptyMemoryStore()
		val := validator.NewBlockValidatorCreator()
		bs := blockstore.New(ms, val)

		open := &Block{Type: OPEN, Link: []byte("ddddddddddddd"), Previous: []byte("ppppppppp"),
			Signature: []byte("a246ce6b1d2b57ac33073127d8f9539fca32fb48481d46d734bf3308796ee18b"), Balance: 1,
			PowNonce: 1, Account: []byte("aaaaaaaaaa"), Representative: []byte("rrrrrrrrrrrrrrr"), Timestamp: 1}
		open.SetHash()

		blk, err := bs.Store(open)

		block := &Block{Type: SEND, Link: []byte("ddddddddddddd"), Previous: open.Hash,
			Signature: []byte("df0d25f706c31d2007ed91da185ac727e5e38bc77f4309bb587e1ff7557ace39"), Balance: 1,
			PowNonce: 1, Account: []byte("aaaaaaaaaa"), Representative: []byte("rrrrrrrrrrrrrrr"), Timestamp: 1,
		}
		block.SetHash()

		blk, err = bs.Store(block)

		Expect(err).To(BeNil())
		Expect(blk).NotTo(BeNil())

		chain, err := bs.GetBlockChain(string(blk.Hash))
		Expect(err).To(BeNil())
		Expect(len(chain)).To(Equal(2))

		Expect(chain[0].Hash).To(Equal([]byte("a246ce6b1d2b57ac33073127d8f9539fca32fb48481d46d734bf3308796ee18b")))
		Expect(chain[1].Hash).To(Equal([]byte("df0d25f706c31d2007ed91da185ac727e5e38bc77f4309bb587e1ff7557ace39")))
	})
})
