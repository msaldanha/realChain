package wallet_test

import (
	"github.com/msaldanha/realChain/tests"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/golang/mock/gomock"
	"github.com/msaldanha/realChain/keyvaluestore"
	"github.com/msaldanha/realChain/address"
	"github.com/msaldanha/realChain/wallet"
	"github.com/msaldanha/realChain/ledger"
)

var ts *ledger.TransactionStore

var _ = Describe("Wallet", func() {
	It("Should send funds if acc has funds to send", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ld := tests.NewMockLedger(mockCtrl)

		wa, firstTx, _ := createWallet(ld)

		tx, err := wa.CreateSendTransaction(string(firstTx.Address), "175jFeuksqWTjChY5L4kAN6pbEtgMSnynM", 300)
		Expect(err).To(BeNil())

		tx, err = ts.Retrieve(string(tx.Hash))
		Expect(err).To(BeNil())
		Expect(tx).NotTo(BeNil())

		Expect(tx.Balance).To(Equal(float64(700)))
	})

	It("Should NOT send funds if acc has not enough funds to send", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ld := tests.NewMockLedger(mockCtrl)

		wa, firstTx, _ := createWallet(ld)

		tx, err := wa.CreateSendTransaction(string(firstTx.Address), "175jFeuksqWTjChY5L4kAN6pbEtgMSnynM", 1300)
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrNotEnoughFunds))
		Expect(tx).To(BeNil())
	})

	It("Should return the list of addresses", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ld := tests.NewMockLedger(mockCtrl)

		wa, _, _ := createWallet(ld)
		wa.CreateAddress()

		addrs, err := wa.GetAddresses()
		Expect(err).To(BeNil())
		Expect(addrs).NotTo(BeNil())
		Expect(len(addrs)).To(Equal(2))
	})

	It("Should return an address", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ld := tests.NewMockLedger(mockCtrl)

		wa, _, _ := createWallet(ld)
		addr2, _ := wa.CreateAddress()

		addr, err := wa.GetAddress([]byte(addr2.Address))
		Expect(err).To(BeNil())
		Expect(addr).To(Equal(addr2))
	})

	It("Should get address' statement", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ld := tests.NewMockLedger(mockCtrl)

		wa, firstTx, _ := createWallet(ld)

		ld.EXPECT().GetAddressStatement(gomock.Any())

		wa.GetAddressStatement(string(firstTx.Address))
	})

	It("Should get the last transaction", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ld := tests.NewMockLedger(mockCtrl)

		wa, firstTx, _ := createWallet(ld)

		ld.EXPECT().GetLastTransaction(string(firstTx.Address))

		wa.GetLastTransaction(string(firstTx.Address))
	})
})

func createFirstTx() (*ledger.Transaction, *address.Address) {
	tx := ledger.NewOpenTransaction()
	addr, _ := address.NewAddressWithKeys()
	tx.Address = []byte(addr.Address)
	tx.Representative = tx.Address
	tx.Balance = 1000
	tx.PubKey = addr.Keys.PublicKey
	tx.SetPow()
	tx.Sign(addr.Keys.ToEcdsaPrivateKey())
	return tx, addr
}

func createWallet(ld ledger.Ledger) (*wallet.Wallet, *ledger.Transaction, *address.Address) {
	ms := keyvaluestore.NewMemoryKeyValueStore()
	as := keyvaluestore.NewMemoryKeyValueStore()
	val := ledger.NewValidatorCreator()
	ts = ledger.NewTransactionStore(ms, val)

	firstTx, addr := createFirstTx()

	as.Put(addr.Address, addr.ToBytes())
	ts.Store(firstTx)

	return wallet.New(ts, as, ld), firstTx, addr
}