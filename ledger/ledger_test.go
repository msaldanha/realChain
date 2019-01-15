package ledger_test

import (
	"github.com/golang/mock/gomock"
	"github.com/msaldanha/realChain/address"
	"github.com/msaldanha/realChain/keyvaluestore"
	"github.com/msaldanha/realChain/ledger"
	"github.com/msaldanha/realChain/tests"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Ledger", func() {

	var ld ledger.Ledger
	var genesisTx *ledger.Transaction
	var genesisAddr *address.Address
	var bs *ledger.TransactionStore
	var ms *keyvaluestore.MemoryKeyValueStore

	BeforeEach(func () {
		ms = keyvaluestore.NewMemoryKeyValueStore()
		val := ledger.NewValidatorCreator()
		bs = ledger.NewTransactionStore(ms, val)

		genesisTx, genesisAddr = tests.CreateGenesisTransaction(1000)

		ld = ledger.NewLocalLedger(bs)
	})

	It("Should initialize the Genesis transaction", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		err := ld.Initialize(genesisTx)
		Expect(err).To(BeNil())

		tx, err := bs.Retrieve(string(genesisTx.Hash))
		Expect(err).To(BeNil())
		Expect(tx).NotTo(BeNil())

		//address.IsValid("13N7kLsMvb5r85hbemba9UFkufXBgiVde")
		//13fxKigmqhBEvQ5qgoJcbMyuL4jnn9guH
	})

	It("Should get a transaction", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		err := ld.Initialize(genesisTx)
		Expect(err).To(BeNil())

		tx, err := ld.GetTransaction(string(genesisTx.Hash))
		Expect(err).To(BeNil())
		Expect(tx).NotTo(BeNil())
		Expect(tx.Hash).To(Equal(genesisTx.Hash))
	})

	It("Should NOT create the Genesis transaction twice", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		err := ld.Initialize(genesisTx)
		Expect(err).To(BeNil())

		err = ld.Initialize(genesisTx)
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrLedgerAlreadyInitialized))
	})

	It("Should send funds if acc has funds to send", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		err := ld.Initialize(genesisTx)
		Expect(err).To(BeNil())

		receiveAddr, err := address.NewAddressWithKeys()
		Expect(err).To(BeNil())

		sendTx, err := tests.CreateSendTransaction(genesisTx, genesisAddr, receiveAddr.Address, 300)
		Expect(err).To(BeNil())

		receiveTx, err := tests.CreateReceiveTransaction(sendTx, 300, receiveAddr, nil)
		Expect(err).To(BeNil())

		err = ld.Register(sendTx, receiveTx)
		Expect(err).To(BeNil())

		tx, err := bs.Retrieve(string(sendTx.Hash))
		Expect(err).To(BeNil())
		Expect(tx).NotTo(BeNil())

		Expect(tx.Balance).To(Equal(float64(700)))
	})

	It("Should NOT send funds to invalid address", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		err := ld.Initialize(genesisTx)
		Expect(err).To(BeNil())

		tx, err := bs.Retrieve(string(genesisTx.Hash))
		Expect(err).To(BeNil())
		Expect(tx).NotTo(BeNil())

		receiveAddr, err := address.NewAddressWithKeys()
		Expect(err).To(BeNil())

		sendTx, err := tests.CreateSendTransaction(genesisTx, genesisAddr, receiveAddr.Address, 300)
		Expect(err).To(BeNil())

		receiveTx, err := tests.CreateReceiveTransaction(sendTx, 300, receiveAddr, nil)
		Expect(err).To(BeNil())

		receiveTx.Address = []byte("xxxxxxxxxxxxx")

		err = ld.Register(sendTx, receiveTx)
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(address.ErrInvalidChecksum))
	})

	It("Should NOT send funds if acc has not enough funds to send", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		err := ld.Initialize(genesisTx)
		Expect(err).To(BeNil())

		tx, err := bs.Retrieve(string(genesisTx.Hash))
		Expect(err).To(BeNil())
		Expect(tx).NotTo(BeNil())

		receiveAddr, err := address.NewAddressWithKeys()
		Expect(err).To(BeNil())

		sendTx, err := tests.CreateSendTransaction(genesisTx, genesisAddr, receiveAddr.Address, 1300)
		Expect(err).To(BeNil())

		receiveTx, err := tests.CreateReceiveTransaction(sendTx, 1300, receiveAddr, nil)
		Expect(err).To(BeNil())

		err = ld.Register(sendTx, receiveTx)
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrNotEnoughFunds))
		Expect(ms.Size()).To(Equal(2))

		Expect(tx.Balance).To(Equal(float64(1000)))
	})

	It("Should receive funds", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		err := ld.Initialize(genesisTx)
		Expect(err).To(BeNil())

		receiveAddr, err := address.NewAddressWithKeys()
		Expect(err).To(BeNil())

		sendTx, err := tests.CreateSendTransaction(genesisTx, genesisAddr, receiveAddr.Address, 400)
		Expect(err).To(BeNil())

		receiveTx, err := tests.CreateReceiveTransaction(sendTx, 400, receiveAddr, nil)
		Expect(err).To(BeNil())

		err = ld.Register(sendTx, receiveTx)
		Expect(err).To(BeNil())

		tx, err := bs.Retrieve(string(sendTx.Hash))
		Expect(err).To(BeNil())
		Expect(tx).NotTo(BeNil())

		Expect(tx.Balance).To(Equal(float64(600)))

		receiveTxFromDb, err := ld.GetLastTransaction(string(receiveTx.Address))
		Expect(err).To(BeNil())
		Expect(receiveTxFromDb).NotTo(BeNil())

		txChain, err := bs.GetTransactionChain(string(receiveTxFromDb.Hash), true)
		Expect(err).To(BeNil())
		Expect(len(txChain)).To(Equal(3))

		Expect(txChain[2].Type).To(Equal(ledger.Transaction_OPEN))
		Expect(txChain[2].Balance).To(Equal(float64(400)))
	})

	It("Should NOT receive funds from tampered transaction", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		err := ld.Initialize(genesisTx)
		Expect(err).To(BeNil())

		receiveAddr, err := address.NewAddressWithKeys()
		Expect(err).To(BeNil())

		sendTx, err := tests.CreateSendTransaction(genesisTx, genesisAddr, receiveAddr.Address, 400)
		Expect(err).To(BeNil())

		receiveTx, err := tests.CreateReceiveTransaction(sendTx, 400, receiveAddr, nil)
		Expect(err).To(BeNil())

		// Tamper with balance
		sendTx.Balance = float64(500)

		err = ld.Register(sendTx, receiveTx)
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrInvalidTransactionHash))

		// Restore balance
		sendTx.Balance = float64(600)
		// Tamper with signature
		sendTx.Signature[0] = sendTx.Signature[0] + 1

		err = ld.Register(sendTx, receiveTx)
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrInvalidTransactionSignature))
	})

	It("Should NOT receive funds from non pending send transaction", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		err := ld.Initialize(genesisTx)
		Expect(err).To(BeNil())

		receiveAddr, err := address.NewAddressWithKeys()
		Expect(err).To(BeNil())

		sendTx, err := tests.CreateSendTransaction(genesisTx, genesisAddr, receiveAddr.Address, 400)
		Expect(err).To(BeNil())

		receiveTx, err := tests.CreateReceiveTransaction(sendTx, 400, receiveAddr, nil)
		Expect(err).To(BeNil())

		err = ld.Register(sendTx, receiveTx)
		Expect(err).To(BeNil())

		tx, err := bs.Retrieve(string(sendTx.Hash))
		Expect(err).To(BeNil())
		Expect(tx).NotTo(BeNil())
		Expect(tx.Balance).To(Equal(float64(600)))

		ms.Put(string(sendTx.Hash), nil)
		ms.Put(string(sendTx.Address), genesisTx.ToBytes())

		newReceiveTx, err := tests.CreateReceiveTransaction(sendTx, 400, receiveAddr, receiveTx)
		Expect(err).To(BeNil())

		err = ld.Register(sendTx, newReceiveTx)
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrSendTransactionIsNotPending))
	})

	It("Should NOT accept transaction when address does not match pub key", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		err := ld.Initialize(genesisTx)
		Expect(err).To(BeNil())

		receiveAddr, err := address.NewAddressWithKeys()
		Expect(err).To(BeNil())

		sendTx, err := tests.CreateSendTransaction(genesisTx, genesisAddr, receiveAddr.Address, 400)
		Expect(err).To(BeNil())

		receiveTx, err := tests.CreateReceiveTransaction(sendTx, 400, receiveAddr, nil)
		Expect(err).To(BeNil())

		sendTx.PubKey[0] = sendTx.PubKey[0] + 1

		err = ld.Register(sendTx, receiveTx)
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrAddressDoesNotMatchPubKey))
	})

	It("Should produce a correct tx chain", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		err := ld.Initialize(genesisTx)
		Expect(err).To(BeNil())

		receiveAddr, err := address.NewAddressWithKeys()
		Expect(err).To(BeNil())

		var prevReceiveTx *ledger.Transaction
		prevSendTx := genesisTx
		for x := 1; x <= 10; x++ {
			prevSendTx, prevReceiveTx = tests.SendFunds(ld, genesisAddr, prevSendTx, prevReceiveTx, receiveAddr, 100)
		}

		txChain, err := bs.GetTransactionChain(string(prevSendTx.Hash), true)
		Expect(err).To(BeNil())
		Expect(len(txChain)).To(Equal(11))
		Expect(txChain[10].Type).To(Equal(ledger.Transaction_SEND))
		Expect(txChain[10].Balance).To(Equal(float64(0)))

		txChain, err = bs.GetTransactionChain(string(prevReceiveTx.Hash), true)
		Expect(err).To(BeNil())
		Expect(len(txChain)).To(Equal(12))
		Expect(txChain[11].Type).To(Equal(ledger.Transaction_RECEIVE))
		Expect(txChain[11].Balance).To(Equal(float64(1000)))
	})

	It("Should produce a correct address statement", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		err := ld.Initialize(genesisTx)
		Expect(err).To(BeNil())

		receiveAddr, err := address.NewAddressWithKeys()
		Expect(err).To(BeNil())

		var prevReceiveTx *ledger.Transaction
		prevSendTx := genesisTx
		for x := 1; x <= 10; x++ {
			prevSendTx, prevReceiveTx = tests.SendFunds(ld, genesisAddr, prevSendTx, prevReceiveTx, receiveAddr, 100)
		}

		tx, _, _ := bs.GetTransaction(string(prevSendTx.Hash))
		txChain, err := ld.GetAddressStatement(string(tx.Address))
		Expect(err).To(BeNil())
		Expect(len(txChain)).To(Equal(11))
		Expect(txChain[0].Type).To(Equal(ledger.Transaction_OPEN))
		Expect(txChain[0].Balance).To(Equal(float64(1000)))
		Expect(txChain[10].Type).To(Equal(ledger.Transaction_SEND))
		Expect(txChain[10].Balance).To(Equal(float64(0)))

		tx, _, _ = bs.GetTransaction(string(prevReceiveTx.Hash))
		txChain, err = ld.GetAddressStatement(string(tx.Address))
		Expect(err).To(BeNil())
		Expect(len(txChain)).To(Equal(10))
		Expect(txChain[0].Type).To(Equal(ledger.Transaction_OPEN))
		Expect(txChain[0].Balance).To(Equal(float64(100)))
		Expect(txChain[9].Type).To(Equal(ledger.Transaction_RECEIVE))
		Expect(txChain[9].Balance).To(Equal(float64(1000)))
	})

	It("Should return correct balance", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		err := ld.Initialize(genesisTx)
		Expect(err).To(BeNil())

		receiveAddr, err := address.NewAddressWithKeys()
		Expect(err).To(BeNil())

		var prevReceiveTx *ledger.Transaction
		prevSendTx := genesisTx
		for x := 1; x <= 2; x++ {
			prevSendTx, prevReceiveTx = tests.SendFunds(ld, genesisAddr, prevSendTx, prevReceiveTx, receiveAddr, 100)
		}

		tx, err := ld.GetLastTransaction(string(prevSendTx.Address))
		Expect(err).To(BeNil())
		Expect(tx).NotTo(BeNil())
		Expect(tx.Balance).To(Equal(float64(800)))

		tx, err = ld.GetLastTransaction(string(prevReceiveTx.Address))
		Expect(err).To(BeNil())
		Expect(tx).NotTo(BeNil())
		Expect(tx.Balance).To(Equal(float64(200)))
	})

	It("Should verify transaction's pow", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		err := ld.Initialize(genesisTx)
		Expect(err).To(BeNil())

		receiveAddr, err := address.NewAddressWithKeys()
		Expect(err).To(BeNil())

		sendTx, err := tests.CreateSendTransaction(genesisTx, genesisAddr, receiveAddr.Address, 350)
		Expect(err).To(BeNil())

		sendTx.Balance = 350

		err = ld.VerifyTransaction(sendTx, true)
		Expect(err).To(Equal(ledger.ErrInvalidTransactionHash))
	})

	It("Should verify transaction's signature", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		err := ld.Initialize(genesisTx)
		Expect(err).To(BeNil())

		receiveAddr, err := address.NewAddressWithKeys()
		Expect(err).To(BeNil())

		sendTx, err := tests.CreateSendTransaction(genesisTx, genesisAddr, receiveAddr.Address, 350)
		Expect(err).To(BeNil())

		sendTx.Signature = sendTx.Hash

		err = ld.VerifyTransaction(sendTx, true)
		Expect(err).To(Equal(ledger.ErrInvalidTransactionSignature))
	})

	It("Should verify if transaction already exists if new", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		err := ld.Initialize(genesisTx)
		Expect(err).To(BeNil())

		receiveAddr, err := address.NewAddressWithKeys()
		Expect(err).To(BeNil())

		sendTx, err := tests.CreateSendTransaction(genesisTx, genesisAddr, receiveAddr.Address, 400)
		Expect(err).To(BeNil())

		receiveTx, err := tests.CreateReceiveTransaction(sendTx, 400, receiveAddr, nil)
		Expect(err).To(BeNil())

		err = ld.Register(sendTx, receiveTx)
		Expect(err).To(BeNil())

		err = ld.VerifyTransaction(sendTx, true)
		Expect(err).To(Equal(ledger.ErrTransactionAlreadyInLedger))
	})

	It("Should verify if transaction exists if not new", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		err := ld.Initialize(genesisTx)
		Expect(err).To(BeNil())

		receiveAddr, err := address.NewAddressWithKeys()
		Expect(err).To(BeNil())

		sendTx, err := tests.CreateSendTransaction(genesisTx, genesisAddr, receiveAddr.Address, 400)
		Expect(err).To(BeNil())

		err = ld.VerifyTransaction(sendTx, false)
		Expect(err).To(Equal(ledger.ErrTransactionNotFound))
	})

	It("Should verify if previous transaction in the chain exists", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		err := ld.Initialize(genesisTx)
		Expect(err).To(BeNil())

		receiveAddr, err := address.NewAddressWithKeys()
		Expect(err).To(BeNil())

		sendTx, err := tests.CreateSendTransaction(genesisTx, genesisAddr, receiveAddr.Address, 400)
		Expect(err).To(BeNil())

		ms.Put(string(genesisTx.Hash), nil)
		ms.Put(string(genesisTx.Address), nil)

		err = ld.VerifyTransaction(sendTx, true)
		Expect(err).To(Equal(ledger.ErrPreviousTransactionNotFound))
	})

	It("Should verify if head transaction in the chain exists", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		err := ld.Initialize(genesisTx)
		Expect(err).To(BeNil())

		receiveAddr, err := address.NewAddressWithKeys()
		Expect(err).To(BeNil())

		sendTx, err := tests.CreateSendTransaction(genesisTx, genesisAddr, receiveAddr.Address, 400)
		Expect(err).To(BeNil())

		ms.Put(string(genesisTx.Address), nil)

		err = ld.VerifyTransaction(sendTx, true)
		Expect(err).To(Equal(ledger.ErrHeadTransactionNotFound))
	})

	It("Should verify if previous transaction in the chain is a head, if inserting a new", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		err := ld.Initialize(genesisTx)
		Expect(err).To(BeNil())

		receiveAddr, err := address.NewAddressWithKeys()
		Expect(err).To(BeNil())

		sendTx, err := tests.CreateSendTransaction(genesisTx, genesisAddr, receiveAddr.Address, 400)
		Expect(err).To(BeNil())

		receiveTx, err := tests.CreateReceiveTransaction(sendTx, 400, receiveAddr, nil)
		Expect(err).To(BeNil())

		err = ld.Register(sendTx, receiveTx)
		Expect(err).To(BeNil())

		sendTx2, err := tests.CreateSendTransaction(genesisTx, genesisAddr, receiveAddr.Address, 100)
		Expect(err).To(BeNil())

		err = ld.VerifyTransaction(sendTx2, true)
		Expect(err).To(Equal(ledger.ErrPreviousTransactionIsNotHead))
	})

	It("Should verify if transaction's balance is negative", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		err := ld.Initialize(genesisTx)
		Expect(err).To(BeNil())

		receiveAddr, err := address.NewAddressWithKeys()
		Expect(err).To(BeNil())

		sendTx, err := tests.CreateSendTransaction(genesisTx, genesisAddr, receiveAddr.Address, -1)
		Expect(err).To(BeNil())

		err = ld.VerifyTransaction(sendTx, true)
		Expect(err).To(Equal(ledger.ErrNotEnoughFunds))
	})

	It("Should verify if acc has enough funds", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		err := ld.Initialize(genesisTx)
		Expect(err).To(BeNil())

		receiveAddr, err := address.NewAddressWithKeys()
		Expect(err).To(BeNil())

		sendTx, err := tests.CreateSendTransaction(genesisTx, genesisAddr, receiveAddr.Address, 1001)
		Expect(err).To(BeNil())

		err = ld.VerifyTransaction(sendTx, true)
		Expect(err).To(Equal(ledger.ErrNotEnoughFunds))
	})

	It("Should verify if acc's open transaction is in the chain", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		err := ld.Initialize(genesisTx)
		Expect(err).To(BeNil())

		receiveAddr, err := address.NewAddressWithKeys()
		Expect(err).To(BeNil())

		sendTx, err := tests.CreateSendTransaction(genesisTx, genesisAddr, receiveAddr.Address, 400)
		Expect(err).To(BeNil())

		receiveTx, err := tests.CreateReceiveTransaction(sendTx, 400, receiveAddr, nil)
		Expect(err).To(BeNil())

		err = ld.Register(sendTx, receiveTx)
		Expect(err).To(BeNil())

		sendTx2, err := tests.CreateSendTransaction(sendTx, genesisAddr, receiveAddr.Address, 100)
		Expect(err).To(BeNil())

		ms.Put(string(genesisTx.Hash), nil)

		err = ld.VerifyTransaction(sendTx2, true)
		Expect(err).To(Equal(ledger.ErrOpenTransactionNotFound))
	})

	It("Should verify if send transaction is of type send", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		err := ld.Initialize(genesisTx)
		Expect(err).To(BeNil())

		receiveAddr, err := address.NewAddressWithKeys()
		Expect(err).To(BeNil())

		sendTx, err := tests.CreateSendTransaction(genesisTx, genesisAddr, receiveAddr.Address, 400)
		Expect(err).To(BeNil())

		receiveTx, err := tests.CreateReceiveTransaction(sendTx, 400, receiveAddr, nil)
		Expect(err).To(BeNil())

		err = ld.Verify(receiveTx, sendTx)
		Expect(err).To(Equal(ledger.ErrInvalidSendTransaction))
	})

	It("Should verify if receive transaction is of type receive", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		err := ld.Initialize(genesisTx)
		Expect(err).To(BeNil())

		receiveAddr, err := address.NewAddressWithKeys()
		Expect(err).To(BeNil())

		sendTx, err := tests.CreateSendTransaction(genesisTx, genesisAddr, receiveAddr.Address, 400)
		Expect(err).To(BeNil())

		err = ld.Verify(sendTx, sendTx)
		Expect(err).To(Equal(ledger.ErrInvalidReceiveTransaction))
	})

	It("Should verify if send and receive transactions are linked", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		err := ld.Initialize(genesisTx)
		Expect(err).To(BeNil())

		receiveAddr, err := address.NewAddressWithKeys()
		Expect(err).To(BeNil())

		sendTx, err := tests.CreateSendTransaction(genesisTx, genesisAddr, receiveAddr.Address, 400)
		Expect(err).To(BeNil())

		sendTx2, err := tests.CreateSendTransaction(genesisTx, genesisAddr, receiveAddr.Address, 600)
		Expect(err).To(BeNil())

		receiveTx, err := tests.CreateReceiveTransaction(sendTx, 400, receiveAddr, nil)
		Expect(err).To(BeNil())

		err = ld.Verify(sendTx2, receiveTx)
		Expect(err).To(Equal(ledger.ErrSendReceiveTransactionsNotLinked))
	})

	It("Should verify if send transaction is pending", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		err := ld.Initialize(genesisTx)
		Expect(err).To(BeNil())

		receiveAddr, err := address.NewAddressWithKeys()
		Expect(err).To(BeNil())

		sendTx, err := tests.CreateSendTransaction(genesisTx, genesisAddr, receiveAddr.Address, 400)
		Expect(err).To(BeNil())

		receiveTx, err := tests.CreateReceiveTransaction(sendTx, 400, receiveAddr, nil)
		Expect(err).To(BeNil())

		err = ld.Register(sendTx, receiveTx)
		Expect(err).To(BeNil())

		_ = ms.Put(string(sendTx.Hash), nil)
		_ = ms.Put(string(sendTx.Address), genesisTx.ToBytes())

		receiveTx2, err := tests.CreateReceiveTransaction(sendTx, 400, receiveAddr, nil)
		Expect(err).To(BeNil())

		err = ld.Verify(sendTx, receiveTx2)
		Expect(err).To(Equal(ledger.ErrSendTransactionIsNotPending))
	})

	It("Should verify if send and receive amount are equal", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		err := ld.Initialize(genesisTx)
		Expect(err).To(BeNil())

		receiveAddr, err := address.NewAddressWithKeys()
		Expect(err).To(BeNil())

		sendTx, err := tests.CreateSendTransaction(genesisTx, genesisAddr, receiveAddr.Address, 400)
		Expect(err).To(BeNil())

		receiveTx, err := tests.CreateReceiveTransaction(sendTx, 300, receiveAddr, nil)
		Expect(err).To(BeNil())

		err = ld.Verify(sendTx, receiveTx)
		Expect(err).To(Equal(ledger.ErrSentAmountDiffersFromReceivedAmount))
	})

	It("Should verify if send and receive are from different accounts", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		err := ld.Initialize(genesisTx)
		Expect(err).To(BeNil())

		sendTx, err := tests.CreateSendTransaction(genesisTx, genesisAddr, genesisAddr.Address, 400)
		Expect(err).To(BeNil())

		receiveTx, err := tests.CreateReceiveTransaction(sendTx, 400, genesisAddr, nil)
		Expect(err).To(BeNil())

		err = ld.Verify(sendTx, receiveTx)
		Expect(err).To(Equal(ledger.ErrSendReceiveTransactionsCantBeSameAddress))
	})
})
