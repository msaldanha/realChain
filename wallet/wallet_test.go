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

	var mockCtrl *gomock.Controller
	var wa *wallet.Wallet
	var firstTx *ledger.Transaction
	var ld *tests.MockLedgerClient

	BeforeEach(func () {
		mockCtrl = gomock.NewController(GinkgoT())
		ld = tests.NewMockLedgerClient(mockCtrl)
		wa, firstTx, _ = createWallet(ld)
	})

	It("Should send funds if acc has funds to send", func() {
		defer mockCtrl.Finish()

		toAddr, _ := wa.CreateAddress()

		expectedToAddr := &ledger.GetLastTransactionRequest{Address: toAddr.Address}
		ld.EXPECT().GetLastTransaction(gomock.Any(), gomock.Eq(expectedToAddr), gomock.Any()).
			Return(&ledger.GetLastTransactionResult{}, nil)

		expectedFromAddr := &ledger.GetLastTransactionRequest{Address: string(firstTx.Address)}
		ld.EXPECT().GetLastTransaction(gomock.Any(), gomock.Eq(expectedFromAddr), gomock.Any()).
			Return(&ledger.GetLastTransactionResult{Tx: firstTx}, nil)

		expectedRegister := &ledger.RegisterResult{}
		ld.EXPECT().Register(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(expectedRegister, nil)

		tx, err := wa.Transfer(string(firstTx.Address), string(toAddr.Address), 300)
		Expect(err).To(BeNil())

		Expect(tx.Balance).To(Equal(float64(300)))
	})

	It("Should NOT send funds if acc has not enough funds to send", func() {
		defer mockCtrl.Finish()

		toAddr, _ := wa.CreateAddress()

		expectedToAddr := &ledger.GetLastTransactionRequest{Address: toAddr.Address}
		ld.EXPECT().GetLastTransaction(gomock.Any(), gomock.Eq(expectedToAddr), gomock.Any()).
			Return(&ledger.GetLastTransactionResult{}, nil)

		expectedFromAddr := &ledger.GetLastTransactionRequest{Address: string(firstTx.Address)}
		ld.EXPECT().GetLastTransaction(gomock.Any(), gomock.Eq(expectedFromAddr), gomock.Any()).
			Return(&ledger.GetLastTransactionResult{Tx: firstTx}, nil)

		tx, err := wa.Transfer(string(firstTx.Address), string(toAddr.Address), 1300)
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrNotEnoughFunds))
		Expect(tx).To(BeNil())
	})

	It("Should return the list of addresses", func() {
		defer mockCtrl.Finish()

		wa.CreateAddress()

		addrs, err := wa.GetAddresses()
		Expect(err).To(BeNil())
		Expect(addrs).NotTo(BeNil())
		Expect(len(addrs)).To(Equal(2))
	})

	It("Should return an address", func() {
		defer mockCtrl.Finish()

		addr2, _ := wa.CreateAddress()

		addr, err := wa.GetAddress([]byte(addr2.Address))
		Expect(err).To(BeNil())
		Expect(addr).To(Equal(addr2))
	})

	It("Should get address' statement", func() {
		defer mockCtrl.Finish()

		ld.EXPECT().GetAddressStatement(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(&ledger.GetAddressStatementResult{Txs:[]*ledger.Transaction{firstTx}}, nil)

		wa.GetAddressStatement(string(firstTx.Address))
	})

	It("Should get the last transaction", func() {
		defer mockCtrl.Finish()

		expectedAddr := &ledger.GetLastTransactionRequest{Address: string(firstTx.Address)}
		ld.EXPECT().GetLastTransaction(gomock.Any(), expectedAddr, gomock.Any()).
			Return(&ledger.GetLastTransactionResult{Tx:firstTx}, nil)

		tx, err := wa.GetLastTransaction(string(firstTx.Address))
		Expect(err).To(BeNil())
		Expect(tx).To(Equal(firstTx))
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

func createWallet(ld ledger.LedgerClient) (*wallet.Wallet, *ledger.Transaction, *address.Address) {
	ms := keyvaluestore.NewMemoryKeyValueStore()
	as := keyvaluestore.NewMemoryKeyValueStore()
	val := ledger.NewValidatorCreator()
	ts = ledger.NewTransactionStore(ms, val)

	firstTx, addr := createFirstTx()

	as.Put(addr.Address, addr.ToBytes())
	ts.Store(firstTx)

	return wallet.New(as, ld), firstTx, addr
}