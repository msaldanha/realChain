package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/msaldanha/realChain/blockstore"
	"github.com/golang/mock/gomock"
	"github.com/msaldanha/realChain/keyvaluestore"
	"github.com/msaldanha/realChain/ledger"
	"github.com/msaldanha/realChain/block"
	"github.com/msaldanha/realChain/address"
)

var _ = Describe("Ledger", func() {
	It("Should create the Genesis block", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		as := keyvaluestore.NewMemoryKeyValueStore()
		val := block.NewBlockValidatorCreator()
		bs := blockstore.New(ms, val)

		ld := ledger.New()
		ld.Use(bs, as)

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
		as := keyvaluestore.NewMemoryKeyValueStore()
		val := block.NewBlockValidatorCreator()
		bs := blockstore.New(ms, val)

		ld := ledger.New()
		ld.Use(bs, as)

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
		as := keyvaluestore.NewMemoryKeyValueStore()
		val := block.NewBlockValidatorCreator()
		bs := blockstore.New(ms, val)

		ld := ledger.New()
		ld.Use(bs, as)

		blk, err := ld.Initialize(1000)
		Expect(err).To(BeNil())

		blk, err = bs.Retrieve(string(blk.Hash))
		Expect(err).To(BeNil())
		Expect(blk).NotTo(BeNil())

		hash, err := ld.Send(string(blk.Account), "175jFeuksqWTjChY5L4kAN6pbEtgMSnynM", 300)
		Expect(err).To(BeNil())

		blk, err = bs.Retrieve(hash)
		Expect(err).To(BeNil())
		Expect(blk).NotTo(BeNil())

		Expect(blk.Balance).To(Equal(float64(700)))
	})

	It("Should NOT send funds to invalid address", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		as := keyvaluestore.NewMemoryKeyValueStore()
		val := block.NewBlockValidatorCreator()
		bs := blockstore.New(ms, val)

		ld := ledger.New()
		ld.Use(bs, as)

		blk, err := ld.Initialize(1000)
		Expect(err).To(BeNil())

		blk, err = bs.Retrieve(string(blk.Hash))
		Expect(err).To(BeNil())
		Expect(blk).NotTo(BeNil())

		hash, err := ld.Send(string(blk.Account), "xxxxxxxxxx", 300)
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(address.ErrInvalidChecksum))
		Expect(hash).To(Equal(""))

		blkSend, err := bs.Retrieve(hash)
		Expect(err).To(BeNil())
		Expect(blkSend).To(BeNil())
	})

	It("Should NOT send funds if acc has not enough funds to send", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		as := keyvaluestore.NewMemoryKeyValueStore()
		val := block.NewBlockValidatorCreator()
		bs := blockstore.New(ms, val)

		ld := ledger.New()
		ld.Use(bs, as)

		blk, err := ld.Initialize(1000)
		Expect(err).To(BeNil())

		blk, err = bs.Retrieve(string(blk.Hash))
		Expect(err).To(BeNil())
		Expect(blk).NotTo(BeNil())

		hash, err := ld.Send(string(blk.Account), "175jFeuksqWTjChY5L4kAN6pbEtgMSnynM", 1200)
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrNotEnoughFunds))
		Expect(hash).To(BeEmpty())
		Expect(ms.Size()).To(Equal(2))

		Expect(blk.Balance).To(Equal(float64(1000)))
	})

	It("Should receive funds", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		as := keyvaluestore.NewMemoryKeyValueStore()
		val := block.NewBlockValidatorCreator()
		bs := blockstore.New(ms, val)

		ld := ledger.New()
		ld.Use(bs, as)

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

		receiveBlk, err := ld.GetLastTransaction(receiveAcc.Address)
		Expect(err).To(BeNil())
		Expect(receiveBlk).NotTo(BeNil())

		blockChain, err := bs.GetBlockChain(string(receiveBlk.Hash), true)
		Expect(err).To(BeNil())
		Expect(len(blockChain)).To(Equal(3))

		Expect(blockChain[2].Type).To(Equal(block.OPEN))
		Expect(blockChain[2].Balance).To(Equal(float64(400)))
	})

	It("Should NOT receive funds from tampered block", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		as := keyvaluestore.NewMemoryKeyValueStore()
		val := block.NewBlockValidatorCreator()
		bs := blockstore.New(ms, val)

		ld := ledger.New()
		ld.Use(bs, as)

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

		sendBlk.Balance = float64(500)

		hash, err = ld.Receive(sendBlk)
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrInvalidTransactionHash))
		Expect(hash).To(Equal(""))

		sendBlk.Balance = float64(600)
		sendBlk.Signature[0] = sendBlk.Signature[0] + 1

		hash, err = ld.Receive(sendBlk)
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrInvalidTransactionSignature))
		Expect(hash).To(Equal(""))
	})

	It("Should NOT receive funds from not pending send transaction", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		as := keyvaluestore.NewMemoryKeyValueStore()
		val := block.NewBlockValidatorCreator()
		bs := blockstore.New(ms, val)

		ld := ledger.New()
		ld.Use(bs, as)

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
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrSendTransactionIsNotPending))
		Expect(hash).To(Equal(""))
	})

	It("Should NOT accept transaction when account does not match pub key", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		as := keyvaluestore.NewMemoryKeyValueStore()
		val := block.NewBlockValidatorCreator()
		bs := blockstore.New(ms, val)

		ld := ledger.New()
		ld.Use(bs, as)

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

		sendBlk.PubKey[0] = sendBlk.PubKey[0] + 1

		hash, err = ld.Receive(sendBlk)
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrAccountDoesNotMatchPubKey))
		Expect(hash).To(Equal(""))
	})

	It("Should produce a correct blockchain", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		as := keyvaluestore.NewMemoryKeyValueStore()
		val := block.NewBlockValidatorCreator()
		bs := blockstore.New(ms, val)

		ld := ledger.New()
		ld.Use(bs, as)

		blk, err := ld.Initialize(1000)
		Expect(err).To(BeNil())

		receiveAcc := createTestAccount()

		ld.AddAccount(receiveAcc)

		blk, err = bs.Retrieve(string(blk.Hash))
		Expect(err).To(BeNil())
		Expect(blk).NotTo(BeNil())

		var sendHash, receiveHash string
		for x := 1; x <= 10; x++ {
			sendHash, receiveHash = sendFunds(ld, bs, blk, receiveAcc.Address, 100)
		}

		blockChain, err := bs.GetBlockChain(sendHash, true)
		dumpBlockChain(blockChain)
		Expect(err).To(BeNil())
		Expect(len(blockChain)).To(Equal(11))
		Expect(blockChain[10].Type).To(Equal(block.SEND))
		Expect(blockChain[10].Balance).To(Equal(float64(0)))

		blockChain, err = bs.GetBlockChain(receiveHash, true)
		dumpBlockChain(blockChain)
		Expect(err).To(BeNil())
		Expect(len(blockChain)).To(Equal(12))
		Expect(blockChain[11].Type).To(Equal(block.RECEIVE))
		Expect(blockChain[11].Balance).To(Equal(float64(1000)))
	})

	It("Should produce a correct account statement", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		as := keyvaluestore.NewMemoryKeyValueStore()
		val := block.NewBlockValidatorCreator()
		bs := blockstore.New(ms, val)

		ld := ledger.New()
		ld.Use(bs, as)

		blk, err := ld.Initialize(1000)
		Expect(err).To(BeNil())

		receiveAcc := createTestAccount()

		ld.AddAccount(receiveAcc)

		blk, err = bs.Retrieve(string(blk.Hash))
		Expect(err).To(BeNil())
		Expect(blk).NotTo(BeNil())

		var sendHash, receiveHash string
		for x := 1; x <= 10; x++ {
			sendHash, receiveHash = sendFunds(ld, bs, blk, receiveAcc.Address, 100)
		}

		blk, _, _ = bs.GetBlock(sendHash)
		blockChain, err := ld.GetAccountStatement(string(blk.Account))
		dumpBlockChain(blockChain)
		Expect(err).To(BeNil())
		Expect(len(blockChain)).To(Equal(11))
		Expect(blockChain[0].Type).To(Equal(block.OPEN))
		Expect(blockChain[0].Balance).To(Equal(float64(1000)))
		Expect(blockChain[10].Type).To(Equal(block.SEND))
		Expect(blockChain[10].Balance).To(Equal(float64(0)))

		blk, _, _ = bs.GetBlock(receiveHash)
		blockChain, err = ld.GetAccountStatement(string(blk.Account))
		dumpBlockChain(blockChain)
		Expect(err).To(BeNil())
		Expect(len(blockChain)).To(Equal(10))
		Expect(blockChain[0].Type).To(Equal(block.OPEN))
		Expect(blockChain[0].Balance).To(Equal(float64(100)))
		Expect(blockChain[9].Type).To(Equal(block.RECEIVE))
		Expect(blockChain[9].Balance).To(Equal(float64(1000)))
	})

	It("Should return correct balance", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		as := keyvaluestore.NewMemoryKeyValueStore()
		val := block.NewBlockValidatorCreator()
		bs := blockstore.New(ms, val)

		ld := ledger.New()
		ld.Use(bs, as)

		blk, err := ld.Initialize(1000)
		Expect(err).To(BeNil())

		receiveAcc := createTestAccount()

		ld.AddAccount(receiveAcc)

		blk, err = bs.Retrieve(string(blk.Hash))
		Expect(err).To(BeNil())
		Expect(blk).NotTo(BeNil())

		for x := 1; x <= 2; x++ {
			sendFunds(ld, bs, blk, receiveAcc.Address, 100)
		}

		blk, err = ld.GetLastTransaction(string(blk.Account))
		Expect(err).To(BeNil())
		Expect(blk).NotTo(BeNil())
		Expect(blk.Balance).To(Equal(float64(800)))

		blk, err = ld.GetLastTransaction(receiveAcc.Address)
		Expect(err).To(BeNil())
		Expect(blk).NotTo(BeNil())
		Expect(blk.Balance).To(Equal(float64(200)))
	})
})

func sendFunds(ld *ledger.Ledger, bs *blockstore.BlockStore, blk *block.Block, receiveAcc string, amount float64) (string, string) {
	sendHash, err := ld.Send(string(blk.Account), receiveAcc, amount)
	Expect(err).To(BeNil())
	sendBlk, err := bs.Retrieve(sendHash)
	Expect(err).To(BeNil())
	Expect(sendBlk).NotTo(BeNil())

	receiveBlk, err := ld.GetLastTransaction(receiveAcc)
	Expect(err).To(BeNil())

	return sendHash, string(receiveBlk.Hash)
}
