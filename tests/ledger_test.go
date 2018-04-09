package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/msaldanha/realChain/blockstore"
	"github.com/golang/mock/gomock"
	"github.com/msaldanha/realChain/keyvaluestore"
	"github.com/msaldanha/realChain/ledger"
	"github.com/msaldanha/realChain/block"
)

var _ = Describe("Ledger", func() {
	It("Should create the Genesis block", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		val := block.NewBlockValidatorCreator()
		bs := blockstore.New(ms, val)

		ld := ledger.New()
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
		val := block.NewBlockValidatorCreator()
		bs := blockstore.New(ms, val)

		ld := ledger.New()
		ld.Use(bs)

		blk, err := ld.Initialize(1000)
		Expect(err).To(BeNil())
		Expect(blk).NotTo(BeNil())

		blk, err = ld.Initialize(1000)
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrLedgerAlreadyInitialized))
	})

	It("Should send funds if acc has funds to send", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		val := block.NewBlockValidatorCreator()
		bs := blockstore.New(ms, val)

		ld := ledger.New()
		ld.Use(bs)

		blk, err := ld.Initialize(1000)
		Expect(err).To(BeNil())

		blk, err = bs.Retrieve(string(blk.Hash))
		Expect(err).To(BeNil())
		Expect(blk).NotTo(BeNil())


		hash, err := ld.Send(string(blk.Account), "xxxxxxxxxx", 300)
		Expect(err).To(BeNil())

		blk, err = bs.Retrieve(hash)
		Expect(err).To(BeNil())
		Expect(blk).NotTo(BeNil())


		Expect(blk.Balance).To(Equal(float64(700)))
	})

	It("Should NOT send funds if acc has not enough funds to send", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		val := block.NewBlockValidatorCreator()
		bs := blockstore.New(ms, val)

		ld := ledger.New()
		ld.Use(bs)

		blk, err := ld.Initialize(1000)
		Expect(err).To(BeNil())

		blk, err = bs.Retrieve(string(blk.Hash))
		Expect(err).To(BeNil())
		Expect(blk).NotTo(BeNil())

		hash, err := ld.Send(string(blk.Account), "xxxxxxxxxx", 1200)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("not enough funds"))
		Expect(hash).To(BeEmpty())
		Expect(ms.Size()).To(Equal(2))

		Expect(blk.Balance).To(Equal(float64(1000)))
	})

	It("Should receive funds", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		val := block.NewBlockValidatorCreator()
		bs := blockstore.New(ms, val)

		ld := ledger.New()
		ld.Use(bs)

		blk, err := ld.Initialize(1000)
		Expect(err).To(BeNil())

		receiveAcc := createTestAccount()

		ld.AddAccount(receiveAcc)

		blk, err = bs.Retrieve(string(blk.Hash))
		Expect(err).To(BeNil())
		Expect(blk).NotTo(BeNil())

		hash, err := ld.Send(string(blk.Account), receiveAcc.Address, 400)
		Expect(err).To(BeNil())

		sendBlk, err := bs.Retrieve(hash)
		Expect(err).To(BeNil())
		Expect(sendBlk).NotTo(BeNil())

		Expect(sendBlk.Balance).To(Equal(float64(600)))

		hash, err = ld.Receive(sendBlk)
		Expect(err).To(BeNil())


		blockChain, err := bs.GetBlockChain(hash)
		Expect(err).To(BeNil())
		Expect(len(blockChain)).To(Equal(3))

		Expect(blockChain[2].Type).To(Equal(block.OPEN))
		Expect(blockChain[2].Balance).To(Equal(float64(400)))
	})

	It("Should produce a correct blockchain", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		val := block.NewBlockValidatorCreator()
		bs := blockstore.New(ms, val)

		ld := ledger.New()
		ld.Use(bs)

		blk, err := ld.Initialize(1000)
		Expect(err).To(BeNil())

		receiveAcc := createTestAccount()

		ld.AddAccount(receiveAcc)

		blk, err = bs.Retrieve(string(blk.Hash))
		Expect(err).To(BeNil())
		Expect(blk).NotTo(BeNil())

		var sendHash, receiveHash string
		for x := 1; x <= 10; x++ {
			sendHash, receiveHash = sendFunds(ld, bs, blk,receiveAcc.Address, 100)
		}

		blockChain, err := bs.GetBlockChain(sendHash)
		dumpBlockChain(blockChain)
		Expect(err).To(BeNil())
		Expect(len(blockChain)).To(Equal(11))
		Expect(blockChain[10].Type).To(Equal(block.SEND))
		Expect(blockChain[10].Balance).To(Equal(float64(0)))

		blockChain, err = bs.GetBlockChain(receiveHash)
		dumpBlockChain(blockChain)
		Expect(err).To(BeNil())
		Expect(len(blockChain)).To(Equal(12))
		Expect(blockChain[11].Type).To(Equal(block.RECEIVE))
		Expect(blockChain[11].Balance).To(Equal(float64(1000)))
	})
})

func sendFunds(ld *ledger.Ledger, bs *blockstore.BlockStore, blk *block.Block, receiveAcc string, amount float64) (string, string) {
	sendHash, err := ld.Send(string(blk.Account), receiveAcc, amount)
	Expect(err).To(BeNil())
	sendBlk, err := bs.Retrieve(sendHash)
	Expect(err).To(BeNil())
	Expect(sendBlk).NotTo(BeNil())

	receiveHash, err := ld.Receive(sendBlk)
	Expect(err).To(BeNil())

	return sendHash, receiveHash
}