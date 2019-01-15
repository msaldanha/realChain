package keyvaluestore_test

import (
	"github.com/golang/mock/gomock"
	"github.com/msaldanha/realChain/address"
	"github.com/msaldanha/realChain/keyvaluestore"
	"github.com/msaldanha/realChain/ledger"
	"github.com/msaldanha/realChain/tests"
	"os"
	"path/filepath"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

		val := ledger.NewValidatorCreator()
		bs := ledger.NewTransactionStore(txStore, val)

		ld := ledger.NewLocalLedger(bs)

		genesisTx, genesisAddr := tests.CreateGenesisTransaction(1000)
		err = ld.Initialize(genesisTx)
		Expect(err).To(BeNil())

		receiveAddr, err := address.NewAddressWithKeys()
		Expect(err).To(BeNil())

		var prevReceiveTx *ledger.Transaction
		prevSendTx := genesisTx
		for x := 1; x <= 10; x++ {
			prevSendTx, prevReceiveTx = tests.SendFunds(ld, genesisAddr, prevSendTx, prevReceiveTx, receiveAddr, 100)
		}

		txChain, err := bs.GetTransactionChain(string(prevSendTx.Hash), true)
		tests.DumpTxChain(txChain)
		Expect(err).To(BeNil())
		Expect(len(txChain)).To(Equal(11))
		Expect(txChain[10].Type).To(Equal(ledger.Transaction_SEND))
		Expect(txChain[10].Balance).To(Equal(float64(0)))

		txChain, err = bs.GetTransactionChain(string(prevReceiveTx.Hash), true)
		tests.DumpTxChain(txChain)
		Expect(err).To(BeNil())
		Expect(len(txChain)).To(Equal(12))
		Expect(txChain[11].Type).To(Equal(ledger.Transaction_RECEIVE))
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

