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

		bklStoreOptions := prepareOptions("BlockChain", filepath.Join(path, "test.db"))
		blkStore := keyvaluestore.NewBoltKeyValueStore()
		err = blkStore.Init(bklStoreOptions)
		Expect(err).To(BeNil())

		accStoreOptions := prepareOptions("Accounts", filepath.Join(path, "accounts-test.db"))
		as := keyvaluestore.NewBoltKeyValueStore()
		err = as.Init(accStoreOptions)
		Expect(err).To(BeNil())

		val := block.NewBlockValidatorCreator()
		bs := blockstore.New(blkStore, val)

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
})

func prepareOptions(bucketName, filepath string) *keyvaluestore.BoltKeyValueStoreOptions {
	options := &keyvaluestore.BoltKeyValueStoreOptions{DbFile: filepath, BucketName: bucketName}
	_, err := os.Stat(options.DbFile)
	if err == nil {
		os.Remove(options.DbFile)
	}
	return options
}

