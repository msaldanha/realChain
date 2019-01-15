package consensus_test

import (
	"github.com/golang/mock/gomock"
	"github.com/msaldanha/realChain/address"
	"github.com/msaldanha/realChain/consensus"
	"github.com/msaldanha/realChain/crypto"
	"github.com/msaldanha/realChain/ledger"
	"github.com/msaldanha/realChain/tests"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Consensus", func() {

	var mockCtrl *gomock.Controller
	var ld *tests.MockLedger
	var con consensus.Consensus
	var sendTx *ledger.Transaction
	var receiveTx *ledger.Transaction

	BeforeEach(func () {
		mockCtrl = gomock.NewController(GinkgoT())
		ld = tests.NewMockLedger(mockCtrl)
		conAddr, err := address.NewAddressWithKeys()
		Expect(err).To(BeNil())
		con = consensus.NewConsensus(ld, conAddr)

		genesisTx, genesisAddr := tests.CreateGenesisTransaction(1000)

		receiveAddr, err := address.NewAddressWithKeys()
		Expect(err).To(BeNil())

		sendTx, err = tests.CreateSendTransaction(genesisTx, genesisAddr, receiveAddr.Address, 300)
		Expect(err).To(BeNil())

		receiveTx, err = tests.CreateReceiveTransaction(sendTx, 300, receiveAddr, nil)
		Expect(err).To(BeNil())
	})

	It("Should vote OK if ledger would accept transactions", func() {
		defer mockCtrl.Finish()

		ld.EXPECT().VerifyTransaction(sendTx, true)
		ld.EXPECT().VerifyTransaction(receiveTx, true)
		ld.EXPECT().Verify(sendTx, receiveTx)

		request := &consensus.VoteRequest{SendTx: sendTx, ReceiveTx: receiveTx}

		vote, err := con.Vote(request)
		Expect(err).To(BeNil())
		Expect(vote.Vote.Ok).To(BeTrue())
	})

	It("Should vote NOK if ledger would NOT accept send transaction", func() {
		defer mockCtrl.Finish()

		ld.EXPECT().VerifyTransaction(sendTx, true).Return(ledger.ErrInvalidSendTransaction)
		ld.EXPECT().VerifyTransaction(receiveTx, true)

		request := &consensus.VoteRequest{SendTx: sendTx, ReceiveTx: receiveTx}

		vote, err := con.Vote(request)
		Expect(err).To(BeNil())
		Expect(vote.Vote.Ok).To(BeFalse())
		Expect(vote.Vote.Reason).To(Equal(ledger.ErrInvalidSendTransaction.Error()))
	})

	It("Should vote NOK if ledger would NOT accept receive transaction", func() {
		defer mockCtrl.Finish()

		ld.EXPECT().VerifyTransaction(receiveTx, true).Return(ledger.ErrInvalidReceiveTransaction)

		request := &consensus.VoteRequest{SendTx: sendTx, ReceiveTx: receiveTx}

		vote, err := con.Vote(request)
		Expect(err).To(BeNil())
		Expect(vote.Vote.Ok).To(BeFalse())
		Expect(vote.Vote.Reason).To(Equal(ledger.ErrInvalidReceiveTransaction.Error()))
	})

	It("Should vote NOK if ledger would NOT accept transaction", func() {
		defer mockCtrl.Finish()

		ld.EXPECT().VerifyTransaction(sendTx, true)
		ld.EXPECT().VerifyTransaction(receiveTx, true)
		ld.EXPECT().Verify(sendTx, receiveTx).Return(ledger.ErrSendReceiveTransactionsNotLinked)

		request := &consensus.VoteRequest{SendTx: sendTx, ReceiveTx: receiveTx}

		vote, err := con.Vote(request)
		Expect(err).To(BeNil())
		Expect(vote.Vote.Ok).To(BeFalse())
		Expect(vote.Vote.Reason).To(Equal(ledger.ErrSendReceiveTransactionsNotLinked.Error()))
	})

	It("Should sign vote", func() {
		defer mockCtrl.Finish()

		ld.EXPECT().VerifyTransaction(sendTx, true)
		ld.EXPECT().VerifyTransaction(receiveTx, true)
		ld.EXPECT().Verify(sendTx, receiveTx).Return(ledger.ErrSendReceiveTransactionsNotLinked)

		request := &consensus.VoteRequest{SendTx: sendTx, ReceiveTx: receiveTx}

		vote, err := con.Vote(request)
		Expect(err).To(BeNil())
		Expect(vote.Vote.Ok).To(BeFalse())
		Expect(vote.Vote.Reason).To(Equal(ledger.ErrSendReceiveTransactionsNotLinked.Error()))
		Expect(vote.Vote.Signature).NotTo(BeNil())
		Expect(crypto.VerifySignature(vote.Vote.Signature, vote.Vote.PubKey, vote.Vote.Hash())).To(BeTrue())
	})

	It("Should accept voted transactions", func() {
		defer mockCtrl.Finish()

		ld.EXPECT().Register(sendTx, receiveTx)

		votes := [1]*consensus.Vote{{Ok: true}}
		request := &consensus.AcceptRequest{SendTx: sendTx, ReceiveTx: receiveTx, Votes: votes[:]}

		vote, err := con.Accept(request)
		Expect(err).To(BeNil())
		Expect(vote).NotTo(BeNil())
	})

	It("Should NOT accept transactions if not total majority", func() {
		defer mockCtrl.Finish()

		votes := [2]*consensus.Vote{{Ok: true}, {Ok: false}}
		request := &consensus.AcceptRequest{SendTx: sendTx, ReceiveTx: receiveTx, Votes: votes[:]}

		vote, err := con.Accept(request)
		Expect(err).To(Equal(consensus.ErrInvalidVotingResult))
		Expect(vote).To(BeNil())
	})

	It("Should NOT accept transactions if there is no votes", func() {
		defer mockCtrl.Finish()

		votes := [0]*consensus.Vote{}
		request := &consensus.AcceptRequest{SendTx: sendTx, ReceiveTx: receiveTx, Votes: votes[:]}

		vote, err := con.Accept(request)
		Expect(err).To(Equal(consensus.ErrInvalidVotingResult))
		Expect(vote).To(BeNil())
	})

	It("Should NOT accept transactions if ledger rejects it", func() {
		defer mockCtrl.Finish()

		ld.EXPECT().Register(sendTx, receiveTx).Return(ledger.ErrSendReceiveTransactionsNotLinked)

		votes := [1]*consensus.Vote{{Ok: true}}
		request := &consensus.AcceptRequest{SendTx: sendTx, ReceiveTx: receiveTx, Votes: votes[:]}

		vote, err := con.Accept(request)
		Expect(err).To(Equal(ledger.ErrSendReceiveTransactionsNotLinked))
		Expect(vote).To(BeNil())
	})
})
