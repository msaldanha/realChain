package server_test

import (
	"github.com/golang/mock/gomock"
	"github.com/msaldanha/realChain/errors"
	"github.com/msaldanha/realChain/address"
	"github.com/msaldanha/realChain/consensus"
	"github.com/msaldanha/realChain/ledger"
	"github.com/msaldanha/realChain/server"
	"github.com/msaldanha/realChain/tests"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Server", func() {

	var mockCtrl *gomock.Controller
	var ld *tests.MockLedger
	var con *tests.MockConsensus
	var dis *tests.MockDiscoverer
	var lis *tests.MockListener
	var conCli *tests.MockConsensusClient
	var srv *server.Server
	var sendTx *ledger.Transaction
	var receiveTx *ledger.Transaction

	BeforeEach(func () {
		mockCtrl = gomock.NewController(GinkgoT())
		ld = tests.NewMockLedger(mockCtrl)
		con = tests.NewMockConsensus(mockCtrl)
		dis = tests.NewMockDiscoverer(mockCtrl)
		lis = tests.NewMockListener(mockCtrl)
		conCli = tests.NewMockConsensusClient(mockCtrl)

		srv = server.New(ld, con, dis, lis)

		genesisTx, genesisAddr := tests.CreateGenesisTransaction(1000)

		receiveAddr, err := address.NewAddressWithKeys()
		Expect(err).To(BeNil())

		sendTx, err = ledger.CreateSendTransaction(genesisTx, genesisAddr, receiveAddr.Address, 300)
		Expect(err).To(BeNil())

		receiveTx, err = ledger.CreateReceiveTransaction(sendTx, 300, receiveAddr, nil)
		Expect(err).To(BeNil())
	})

	It("Should register transactions in the ledger", func() {
		defer mockCtrl.Finish()

		ld.EXPECT().Register(sendTx, receiveTx)
		ld.EXPECT().Verify(sendTx, receiveTx)

		conCli.EXPECT().Vote(gomock.Any(), gomock.Any(), gomock.Any()).Return(&consensus.VoteResult{Vote:&consensus.Vote{Ok:true}}, nil)
		conCli.EXPECT().Accept(gomock.Any(), gomock.Any(), gomock.Any())

		dis.EXPECT().Peers().Return([]consensus.ConsensusClient{conCli}, nil)

		request := &ledger.RegisterRequest{SendTx: sendTx, ReceiveTx: receiveTx}

		result, err := srv.Register(nil, request)

		Expect(err).To(BeNil())
		Expect(result).NotTo(BeNil())
	})

	It("Should handle error from ledger register transactions", func() {
		defer mockCtrl.Finish()

		ld.EXPECT().Register(sendTx, receiveTx).Return(ledger.ErrInvalidReceiveTransaction)
		ld.EXPECT().Verify(sendTx, receiveTx)

		conCli.EXPECT().Vote(gomock.Any(), gomock.Any(), gomock.Any()).Return(&consensus.VoteResult{Vote:&consensus.Vote{Ok:true}}, nil)

		dis.EXPECT().Peers().Return([]consensus.ConsensusClient{conCli}, nil)

		request := &ledger.RegisterRequest{SendTx: sendTx, ReceiveTx: receiveTx}

		result, err := srv.Register(nil, request)

		Expect(result).To(BeNil())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrInvalidReceiveTransaction))
	})

	It("Should return error if ledger.Verify returns error", func() {
		defer mockCtrl.Finish()

		ld.EXPECT().Verify(sendTx, receiveTx).Return(ledger.ErrPreviousTransactionIsNotHead)

		request := &ledger.RegisterRequest{SendTx: sendTx, ReceiveTx: receiveTx}

		result, err := srv.Register(nil, request)

		Expect(result).To(BeNil())
		Expect(err).To(Equal(ledger.ErrPreviousTransactionIsNotHead))
	})

	It("Should return error if there are no peers", func() {
		defer mockCtrl.Finish()

		ld.EXPECT().Verify(sendTx, receiveTx)

		peers := [0]consensus.ConsensusClient{}
		dis.EXPECT().Peers().Return(peers[:], nil)

		request := &ledger.RegisterRequest{SendTx: sendTx, ReceiveTx: receiveTx}

		result, err := srv.Register(nil, request)

		Expect(result).To(BeNil())
		Expect(err).To(Equal(server.ErrNoPeersForVoting))
	})

	It("Should return error if conflict resolution fails due to error", func() {
		defer mockCtrl.Finish()

		ld.EXPECT().Verify(sendTx, receiveTx)

		expectedError := errors.Error("some error")
		conCli.EXPECT().Vote(gomock.Any(), gomock.Any()).Return(nil, expectedError)
		peers := [1]consensus.ConsensusClient{conCli}
		dis.EXPECT().Peers().Return(peers[:], nil)

		request := &ledger.RegisterRequest{SendTx: sendTx, ReceiveTx: receiveTx}

		result, err := srv.Register(nil, request)

		Expect(result).To(BeNil())
		Expect(err).To(Equal(expectedError))
	})

	It("Should return ErrDeclinedByVoting error if conflict resolution does not approves the transaction", func() {
		defer mockCtrl.Finish()

		ld.EXPECT().Verify(sendTx, receiveTx)

		conCli.EXPECT().Vote(gomock.Any(), gomock.Any()).Return(&consensus.VoteResult{Vote: &consensus.Vote{Ok: false}}, nil)
		peers := [1]consensus.ConsensusClient{conCli}
		dis.EXPECT().Peers().Return(peers[:], nil)

		request := &ledger.RegisterRequest{SendTx: sendTx, ReceiveTx: receiveTx}

		result, err := srv.Register(nil, request)

		Expect(result).To(BeNil())
		Expect(err).To(Equal(server.ErrDeclinedByVoting))
	})

	It("Should call get last transaction from ledger", func() {
		defer mockCtrl.Finish()

		request := &ledger.GetLastTransactionRequest{Address: "xxxxxx"}
		ld.EXPECT().GetLastTransaction(request.Address).Return(sendTx, nil)

		result, err := srv.GetLastTransaction(nil, request)

		Expect(result).NotTo(BeNil())
		Expect(err).To(BeNil())
		Expect(result.Tx).To(Equal(sendTx))
	})

	It("Should return error if call to get last transaction returns error", func() {
		defer mockCtrl.Finish()

		request := &ledger.GetLastTransactionRequest{Address: "xxxxxx"}
		ld.EXPECT().GetLastTransaction(request.Address).Return(nil, ledger.ErrTransactionNotFound)

		_, err := srv.GetLastTransaction(nil, request)

		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrTransactionNotFound))
	})

	It("Should call get transaction from ledger", func() {
		defer mockCtrl.Finish()

		request := &ledger.GetTransactionRequest{Hash: "xxxxxx"}
		ld.EXPECT().GetTransaction(request.Hash).Return(sendTx, nil)

		result, err := srv.GetTransaction(nil, request)

		Expect(result).NotTo(BeNil())
		Expect(err).To(BeNil())
		Expect(result.Tx).To(Equal(sendTx))
	})

	It("Should return error if call to get transaction returns error", func() {
		defer mockCtrl.Finish()

		request := &ledger.GetTransactionRequest{Hash: "xxxxxx"}
		ld.EXPECT().GetTransaction(request.Hash).Return(nil, ledger.ErrTransactionNotFound)

		_, err := srv.GetTransaction(nil, request)

		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrTransactionNotFound))
	})

	It("Should call verify transaction from ledger", func() {
		defer mockCtrl.Finish()

		request := &ledger.VerifyTransactionRequest{Tx: sendTx}
		ld.EXPECT().VerifyTransaction(request.Tx, true).Return(nil)

		result, err := srv.VerifyTransaction(nil, request)

		Expect(result).NotTo(BeNil())
		Expect(err).To(BeNil())
	})

	It("Should return error if call to verify transaction returns error", func() {
		defer mockCtrl.Finish()

		request := &ledger.VerifyTransactionRequest{Tx: sendTx}
		ld.EXPECT().VerifyTransaction(request.Tx, true).Return(ledger.ErrInvalidTransactionHash)

		result, err := srv.VerifyTransaction(nil, request)

		Expect(result).To(BeNil())
		Expect(err).To(Equal(ledger.ErrInvalidTransactionHash))
	})

	It("Should call verify transactions from ledger", func() {
		defer mockCtrl.Finish()

		request := &ledger.VerifyRequest{SendTx: sendTx, ReceiveTx: receiveTx}
		ld.EXPECT().Verify(request.SendTx, request.ReceiveTx).Return(nil)

		result, err := srv.Verify(nil, request)

		Expect(result).NotTo(BeNil())
		Expect(err).To(BeNil())
	})

	It("Should return error if call to verify transactions returns error", func() {
		defer mockCtrl.Finish()

		request := &ledger.VerifyRequest{SendTx: sendTx, ReceiveTx: receiveTx}
		ld.EXPECT().Verify(request.SendTx, request.ReceiveTx).Return(ledger.ErrSendReceiveTransactionsNotLinked)

		result, err := srv.Verify(nil, request)

		Expect(result).To(BeNil())
		Expect(err).To(Equal(ledger.ErrSendReceiveTransactionsNotLinked))
	})

	It("Should call Vote from consensus", func() {
		defer mockCtrl.Finish()

		someErr := errors.Error("some error")
		con.EXPECT().Vote(gomock.Any()).Return(&consensus.VoteResult{Vote: &consensus.Vote{Ok: false, Reason: someErr.Error()}}, nil)

		request := &consensus.VoteRequest{SendTx: sendTx, ReceiveTx: receiveTx}
		result, err := srv.Vote(nil, request)

		Expect(err).To(BeNil())
		Expect(result).NotTo(BeNil())
		Expect(result.Vote.Ok).To(BeFalse())
		Expect(result.Vote.Reason).To(Equal(someErr.Error()))
	})

	It("Should return error if call to Vote returns error", func() {
		defer mockCtrl.Finish()

		someErr := errors.Error("some error")
		con.EXPECT().Vote(gomock.Any()).Return(nil, someErr)

		request := &consensus.VoteRequest{SendTx: sendTx, ReceiveTx: receiveTx}
		result, err := srv.Vote(nil, request)

		Expect(result).To(BeNil())
		Expect(err).To(Equal(someErr))
	})
})
