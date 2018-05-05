package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/msaldanha/realChain/transactionstore"
	"github.com/golang/mock/gomock"
	"github.com/msaldanha/realChain/keyvaluestore"
	"github.com/msaldanha/realChain/ledger"
	"github.com/msaldanha/realChain/transaction"
	"os"
	"path/filepath"
)

var _ = Describe("BoltKeyValueStore", func() {
	It("Should save a correct transaction chain", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		path, err := os.Getwd()

		bklStoreOptions := prepareOptions("TxChain", filepath.Join(path, "test.db"))
		txStore := keyvaluestore.NewBoltKeyValueStore()
		err = txStore.Init(bklStoreOptions)
		Expect(err).To(BeNil())

		accStoreOptions := prepareOptions("Addresses", filepath.Join(path, "addresses-test.db"))
		as := keyvaluestore.NewBoltKeyValueStore()
		err = as.Init(accStoreOptions)
		Expect(err).To(BeNil())

		val := transaction.NewValidatorCreator()
		bs := transactionstore.New(txStore, val)

		ld := ledger.NewLocalLedger(bs, as)

		tx, addr, err := ld.Initialize(1000)
		Expect(err).To(BeNil())

		receiveAcc := createTestAddress()

		ld.AddAddress(receiveAcc)

		tx, err = bs.Retrieve(string(tx.Hash))
		Expect(err).To(BeNil())
		Expect(tx).NotTo(BeNil())

		var sendHash, receiveHash string
		prevTx := tx
		for x := 1; x <= 10; x++ {
			prevTx, sendHash, receiveHash = sendFunds(ld, addr, prevTx, receiveAcc.Address, 100)
		}

		txChain, err := bs.GetTransactionChain(sendHash, true)
		dumpTxChain(txChain)
		Expect(err).To(BeNil())
		Expect(len(txChain)).To(Equal(11))
		Expect(txChain[10].Type).To(Equal(transaction.SEND))
		Expect(txChain[10].Balance).To(Equal(float64(0)))

		txChain, err = bs.GetTransactionChain(receiveHash, true)
		dumpTxChain(txChain)
		Expect(err).To(BeNil())
		Expect(len(txChain)).To(Equal(12))
		Expect(txChain[11].Type).To(Equal(transaction.RECEIVE))
		Expect(txChain[11].Balance).To(Equal(float64(1000)))
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

