package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/msaldanha/realChain/blockstore"
	"github.com/golang/mock/gomock"
	"github.com/msaldanha/realChain/keyvaluestore"
	"github.com/msaldanha/realChain/ledger"
	"github.com/msaldanha/realChain/block"
	"os"
	"path/filepath"
)

var _ = Describe("BoltKeyValueStore", func() {
	It("Should save a correct blockchain", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		path, err := os.Getwd()

		options := &keyvaluestore.BoltKeyValueStoreOptions{DbFile:filepath.Join(path, "test.db")}
		_, err = os.Stat(options.DbFile)
		if err == nil {
			os.Remove(options.DbFile)
		}

		ms := keyvaluestore.NewBoltKeyValueStore()
		err = ms.Init(options)
		Expect(err).To(BeNil())

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
			sendHash, receiveHash = sendFunds(ld, bs, blk, receiveAcc.Address, 100)
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

