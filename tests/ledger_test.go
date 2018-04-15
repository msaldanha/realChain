package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/msaldanha/realChain/transactionstore"
	"github.com/golang/mock/gomock"
	"github.com/msaldanha/realChain/keyvaluestore"
	"github.com/msaldanha/realChain/ledger"
	"github.com/msaldanha/realChain/transaction"
	"github.com/msaldanha/realChain/address"
)

var _ = Describe("Ledger", func() {
	It("Should create the Genesis transaction", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		as := keyvaluestore.NewMemoryKeyValueStore()
		val := transaction.NewValidatorCreator()
		bs := transactionstore.New(ms, val)

		ld := ledger.New()
		ld.Use(bs, as)

		tx, err := ld.Initialize(1000)
		Expect(err).To(BeNil())

		tx, err = bs.Retrieve(string(tx.Hash))
		Expect(err).To(BeNil())
		Expect(tx).NotTo(BeNil())
	})

	It("Should NOT create the Genesis transaction twice", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		as := keyvaluestore.NewMemoryKeyValueStore()
		val := transaction.NewValidatorCreator()
		bs := transactionstore.New(ms, val)

		ld := ledger.New()
		ld.Use(bs, as)

		tx, err := ld.Initialize(1000)
		Expect(err).To(BeNil())
		Expect(tx).NotTo(BeNil())

		tx, err = ld.Initialize(1000)
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrLedgerAlreadyInitialized))
	})

	It("Should send funds if acc has funds to send", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		as := keyvaluestore.NewMemoryKeyValueStore()
		val := transaction.NewValidatorCreator()
		bs := transactionstore.New(ms, val)

		ld := ledger.New()
		ld.Use(bs, as)

		tx, err := ld.Initialize(1000)
		Expect(err).To(BeNil())

		tx, err = bs.Retrieve(string(tx.Hash))
		Expect(err).To(BeNil())
		Expect(tx).NotTo(BeNil())

		tx, err = ld.CreateSendTransaction(string(tx.Account), "175jFeuksqWTjChY5L4kAN6pbEtgMSnynM", 300)
		Expect(err).To(BeNil())

		err = ld.HandleTransaction(tx)
		Expect(err).To(BeNil())

		tx, err = bs.Retrieve(string(tx.Hash))
		Expect(err).To(BeNil())
		Expect(tx).NotTo(BeNil())

		Expect(tx.Balance).To(Equal(float64(700)))
	})

	It("Should NOT send funds to invalid address", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		as := keyvaluestore.NewMemoryKeyValueStore()
		val := transaction.NewValidatorCreator()
		bs := transactionstore.New(ms, val)

		ld := ledger.New()
		ld.Use(bs, as)

		tx, err := ld.Initialize(1000)
		Expect(err).To(BeNil())

		tx, err = bs.Retrieve(string(tx.Hash))
		Expect(err).To(BeNil())
		Expect(tx).NotTo(BeNil())

		tx, err = ld.CreateSendTransaction(string(tx.Account), "xxxxxxxxxx", 300)
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(address.ErrInvalidChecksum))
		Expect(tx).To(BeNil())
	})

	It("Should NOT send funds if acc has not enough funds to send", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		as := keyvaluestore.NewMemoryKeyValueStore()
		val := transaction.NewValidatorCreator()
		bs := transactionstore.New(ms, val)

		ld := ledger.New()
		ld.Use(bs, as)

		tx, err := ld.Initialize(1000)
		Expect(err).To(BeNil())

		tx, err = bs.Retrieve(string(tx.Hash))
		Expect(err).To(BeNil())
		Expect(tx).NotTo(BeNil())

		tx1, err := ld.CreateSendTransaction(string(tx.Account), "175jFeuksqWTjChY5L4kAN6pbEtgMSnynM", 1200)
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrNotEnoughFunds))
		Expect(tx1).To(BeNil())
		Expect(ms.Size()).To(Equal(2))

		Expect(tx.Balance).To(Equal(float64(1000)))
	})

	It("Should receive funds", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		as := keyvaluestore.NewMemoryKeyValueStore()
		val := transaction.NewValidatorCreator()
		bs := transactionstore.New(ms, val)

		ld := ledger.New()
		ld.Use(bs, as)

		tx, err := ld.Initialize(1000)
		Expect(err).To(BeNil())

		receiveAcc := createTestAccount()

		ld.AddAccount(receiveAcc)

		tx, err = bs.Retrieve(string(tx.Hash))
		Expect(err).To(BeNil())
		Expect(tx).NotTo(BeNil())

		tx, err = ld.CreateSendTransaction(string(tx.Account), receiveAcc.Address, 400)
		Expect(err).To(BeNil())

		err = ld.HandleTransaction(tx)
		Expect(err).To(BeNil())

		sendTx, err := bs.Retrieve(string(tx.Hash))
		Expect(err).To(BeNil())
		Expect(sendTx).NotTo(BeNil())

		Expect(sendTx.Balance).To(Equal(float64(600)))

		receiveTx, err := ld.GetLastTransaction(receiveAcc.Address)
		Expect(err).To(BeNil())
		Expect(receiveTx).NotTo(BeNil())

		txChain, err := bs.GetTransactionChain(string(receiveTx.Hash), true)
		Expect(err).To(BeNil())
		Expect(len(txChain)).To(Equal(3))

		Expect(txChain[2].Type).To(Equal(transaction.OPEN))
		Expect(txChain[2].Balance).To(Equal(float64(400)))
	})

	It("Should NOT receive funds from tampered transaction", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		as := keyvaluestore.NewMemoryKeyValueStore()
		val := transaction.NewValidatorCreator()
		bs := transactionstore.New(ms, val)

		ld := ledger.New()
		ld.Use(bs, as)

		tx, err := ld.Initialize(1000)
		Expect(err).To(BeNil())

		receiveAcc := createTestAccount()

		ld.AddAccount(receiveAcc)

		tx, err = bs.Retrieve(string(tx.Hash))
		Expect(err).To(BeNil())
		Expect(tx).NotTo(BeNil())

		tx, err = ld.CreateSendTransaction(string(tx.Account), receiveAcc.Address, 400)
		Expect(err).To(BeNil())

		err = ld.HandleTransaction(tx)
		Expect(err).To(BeNil())

		sendTx, err := bs.Retrieve(string(tx.Hash))
		Expect(err).To(BeNil())
		Expect(sendTx).NotTo(BeNil())

		Expect(sendTx.Balance).To(Equal(float64(600)))

		sendTx.Balance = float64(500)

		hash, err := ld.Receive(sendTx)
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrInvalidTransactionHash))
		Expect(hash).To(Equal(""))

		sendTx.Balance = float64(600)
		sendTx.Signature[0] = sendTx.Signature[0] + 1

		hash, err = ld.Receive(sendTx)
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrInvalidTransactionSignature))
		Expect(hash).To(Equal(""))
	})

	It("Should NOT receive funds from not pending send transaction", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		as := keyvaluestore.NewMemoryKeyValueStore()
		val := transaction.NewValidatorCreator()
		bs := transactionstore.New(ms, val)

		ld := ledger.New()
		ld.Use(bs, as)

		tx, err := ld.Initialize(1000)
		Expect(err).To(BeNil())

		receiveAcc := createTestAccount()

		ld.AddAccount(receiveAcc)

		tx, err = bs.Retrieve(string(tx.Hash))
		Expect(err).To(BeNil())
		Expect(tx).NotTo(BeNil())

		tx, err = ld.CreateSendTransaction(string(tx.Account), receiveAcc.Address, 400)
		Expect(err).To(BeNil())

		err = ld.HandleTransaction(tx)
		Expect(err).To(BeNil())

		sendTx, err := bs.Retrieve(string(tx.Hash))
		Expect(err).To(BeNil())
		Expect(sendTx).NotTo(BeNil())

		Expect(sendTx.Balance).To(Equal(float64(600)))

		hash, err := ld.Receive(sendTx)
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrSendTransactionIsNotPending))
		Expect(hash).To(Equal(""))
	})

	It("Should NOT accept transaction when account does not match pub key", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		as := keyvaluestore.NewMemoryKeyValueStore()
		val := transaction.NewValidatorCreator()
		bs := transactionstore.New(ms, val)

		ld := ledger.New()
		ld.Use(bs, as)

		tx, err := ld.Initialize(1000)
		Expect(err).To(BeNil())

		receiveAcc := createTestAccount()

		ld.AddAccount(receiveAcc)

		tx, err = bs.Retrieve(string(tx.Hash))
		Expect(err).To(BeNil())
		Expect(tx).NotTo(BeNil())

		tx, err = ld.CreateSendTransaction(string(tx.Account), receiveAcc.Address, 400)
		Expect(err).To(BeNil())

		err = ld.HandleTransaction(tx)
		Expect(err).To(BeNil())

		sendTx, err := bs.Retrieve(string(tx.Hash))
		Expect(err).To(BeNil())
		Expect(sendTx).NotTo(BeNil())

		Expect(sendTx.Balance).To(Equal(float64(600)))

		sendTx.PubKey[0] = sendTx.PubKey[0] + 1

		hash, err := ld.Receive(sendTx)
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrAccountDoesNotMatchPubKey))
		Expect(hash).To(Equal(""))
	})

	It("Should produce a correct tx chain", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		as := keyvaluestore.NewMemoryKeyValueStore()
		val := transaction.NewValidatorCreator()
		bs := transactionstore.New(ms, val)

		ld := ledger.New()
		ld.Use(bs, as)

		tx, err := ld.Initialize(1000)
		Expect(err).To(BeNil())

		receiveAcc := createTestAccount()

		ld.AddAccount(receiveAcc)

		tx, err = bs.Retrieve(string(tx.Hash))
		Expect(err).To(BeNil())
		Expect(tx).NotTo(BeNil())

		var sendHash, receiveHash string
		for x := 1; x <= 10; x++ {
			sendHash, receiveHash = sendFunds(ld, bs, tx, receiveAcc.Address, 100)
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

	It("Should produce a correct account statement", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		as := keyvaluestore.NewMemoryKeyValueStore()
		val := transaction.NewValidatorCreator()
		bs := transactionstore.New(ms, val)

		ld := ledger.New()
		ld.Use(bs, as)

		tx, err := ld.Initialize(1000)
		Expect(err).To(BeNil())

		receiveAcc := createTestAccount()

		ld.AddAccount(receiveAcc)

		tx, err = bs.Retrieve(string(tx.Hash))
		Expect(err).To(BeNil())
		Expect(tx).NotTo(BeNil())

		var sendHash, receiveHash string
		for x := 1; x <= 10; x++ {
			sendHash, receiveHash = sendFunds(ld, bs, tx, receiveAcc.Address, 100)
		}

		tx, _, _ = bs.GetTransaction(sendHash)
		txChain, err := ld.GetAccountStatement(string(tx.Account))
		dumpTxChain(txChain)
		Expect(err).To(BeNil())
		Expect(len(txChain)).To(Equal(11))
		Expect(txChain[0].Type).To(Equal(transaction.OPEN))
		Expect(txChain[0].Balance).To(Equal(float64(1000)))
		Expect(txChain[10].Type).To(Equal(transaction.SEND))
		Expect(txChain[10].Balance).To(Equal(float64(0)))

		tx, _, _ = bs.GetTransaction(receiveHash)
		txChain, err = ld.GetAccountStatement(string(tx.Account))
		dumpTxChain(txChain)
		Expect(err).To(BeNil())
		Expect(len(txChain)).To(Equal(10))
		Expect(txChain[0].Type).To(Equal(transaction.OPEN))
		Expect(txChain[0].Balance).To(Equal(float64(100)))
		Expect(txChain[9].Type).To(Equal(transaction.RECEIVE))
		Expect(txChain[9].Balance).To(Equal(float64(1000)))
	})

	It("Should return correct balance", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := keyvaluestore.NewMemoryKeyValueStore()
		as := keyvaluestore.NewMemoryKeyValueStore()
		val := transaction.NewValidatorCreator()
		bs := transactionstore.New(ms, val)

		ld := ledger.New()
		ld.Use(bs, as)

		tx, err := ld.Initialize(1000)
		Expect(err).To(BeNil())

		receiveAcc := createTestAccount()

		ld.AddAccount(receiveAcc)

		tx, err = bs.Retrieve(string(tx.Hash))
		Expect(err).To(BeNil())
		Expect(tx).NotTo(BeNil())

		for x := 1; x <= 2; x++ {
			sendFunds(ld, bs, tx, receiveAcc.Address, 100)
		}

		tx, err = ld.GetLastTransaction(string(tx.Account))
		Expect(err).To(BeNil())
		Expect(tx).NotTo(BeNil())
		Expect(tx.Balance).To(Equal(float64(800)))

		tx, err = ld.GetLastTransaction(receiveAcc.Address)
		Expect(err).To(BeNil())
		Expect(tx).NotTo(BeNil())
		Expect(tx.Balance).To(Equal(float64(200)))
	})
})

func sendFunds(ld *ledger.Ledger, bs *transactionstore.TransactionStore, tx *transaction.Transaction, receiveAcc string, amount float64) (string, string) {
	sendTx, err := ld.CreateSendTransaction(string(tx.Account), receiveAcc, amount)
	Expect(err).To(BeNil())

	err = ld.HandleTransaction(sendTx)
	Expect(err).To(BeNil())

	Expect(sendTx).NotTo(BeNil())

	receiveTx, err := ld.GetLastTransaction(receiveAcc)
	Expect(err).To(BeNil())

	return string(sendTx.Hash), string(receiveTx.Hash)
}
