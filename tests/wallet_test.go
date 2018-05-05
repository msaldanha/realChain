package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/msaldanha/realChain/transactionstore"
	"github.com/golang/mock/gomock"
	"github.com/msaldanha/realChain/keyvaluestore"
	"github.com/msaldanha/realChain/transaction"
	"github.com/msaldanha/realChain/address"
	"github.com/msaldanha/realChain/wallet"
)

var _ = Describe("Wallet", func() {
	It("Should send funds if acc has funds to send", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		as := keyvaluestore.NewMemoryKeyValueStore()
		val := transaction.NewValidatorCreator()
		ts := transactionstore.New(ms, val)

		firstTx, addr := createFirstTx()

		as.Put(addr.Address, addr.ToBytes())
		ts.Store(firstTx)

		ld := wallet.New(ts, as)

		tx, err := ld.CreateSendTransaction(string(firstTx.Address), "175jFeuksqWTjChY5L4kAN6pbEtgMSnynM", 300)
		Expect(err).To(BeNil())

		tx, err = ts.Retrieve(string(tx.Hash))
		Expect(err).To(BeNil())
		Expect(tx).NotTo(BeNil())

		Expect(tx.Balance).To(Equal(float64(700)))
	})
})

func createFirstTx() (*transaction.Transaction, *address.Address) {
	tx := transaction.NewOpenTransaction()
	addr, _ := address.NewAddressWithKeys()
	tx.Address = []byte(addr.Address)
	tx.Representative = tx.Address
	tx.Balance = 1000
	tx.PubKey = addr.Keys.PublicKey
	tx.SetPow()
	tx.Sign(addr.Keys.ToEcdsaPrivateKey())
	return tx, addr
}