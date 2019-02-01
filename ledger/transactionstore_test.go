package ledger_test

import (
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/proto"
	"github.com/msaldanha/realChain/ledger"
	"github.com/msaldanha/realChain/tests"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("TransactionStore", func() {

	It("Should not accept empty/partially filled transaction", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := tests.CreateNonEmptyMemoryStore()
		val := ledger.NewValidatorCreator()
		bs := ledger.NewTransactionStore(ms, val)
		tx := &ledger.Transaction{}
		tx, err := bs.Store(tx)
		Expect(tx).To(BeNil())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrInvalidTransactionType))

		tx = &ledger.Transaction{Type: ledger.Transaction_SEND, Link: "dddddddddddddd", Previous: "eeeeeeeeee", Signature: "ffffffff", Balance: 1,
			PowNonce: 1, Address: "aaaaaaaaaa", Timestamp: time.Now().Unix(), PubKey: "bbbbbbbb"}
		tx1, err := bs.Store(tx)
		Expect(tx1).To(BeNil())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrTransactionHashCantBeEmpty))

		tx.SetHash()

		dest := &ledger.Transaction{Type: ledger.Transaction_OPEN, Link: "dddddddddddddd", Previous: "eeeeeeeeee", Signature: "ffffffff", Balance: 1,
			PowNonce: 1, Address: "aaaaaaaaaa", Timestamp: time.Now().Unix()}
		ms.Put("dddddddddddddd", dest.ToBytes())

		tx, err = bs.Store(tx)
		Expect(err).To(BeNil())
		Expect(tx).NotTo(BeNil())
	})

	It("Should accept properly filled transaction", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := tests.CreateNonEmptyMemoryStore()
		val := ledger.NewValidatorCreator()
		bs := ledger.NewTransactionStore(ms, val)

		tx := &ledger.Transaction{Type: ledger.Transaction_SEND, Link: "dddddddddddddd", Previous: "eeeeeeeeee",
			Signature: "a246ce6b1d2b57ac33073127d8f9539fca32fb48481d46d734bf3308796ee18b", Balance: 1,
			PowNonce: 1, Address: "aaaaaaaaaa", Timestamp: 1, PubKey: "bbbbbbbb"}
		tx.SetHash()

		dest := &ledger.Transaction{Type: ledger.Transaction_OPEN, Link: "dddddddddddddd", Previous: "eeeeeeeeee", Signature: "ffffffff", Balance: 1,
			PowNonce: 1, Address: "aaaaaaaaaa", Timestamp: 1}
		ms.Put("dddddddddddddd", dest.ToBytes())

		tx, err := bs.Store(tx)
		Expect(err).To(BeNil())
		Expect(tx).NotTo(BeNil())

		txFromKeyStore, _, _ := ms.Get(string(tx.Hash))
		Expect(proto.Equal(ledger.NewTransactionFromBytes(txFromKeyStore), tx)).To(BeTrue())
	})

	It("Should extract the chain for the transaction", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := tests.CreateNonEmptyMemoryStore()
		val := ledger.NewValidatorCreator()
		bs := ledger.NewTransactionStore(ms, val)

		open := &ledger.Transaction{Type: ledger.Transaction_OPEN, Link: "dddddddddddddd", Previous: "eeeeeeeeee",
			Signature: "a246ce6b1d2b57ac33073127d8f9539fca32fb48481d46d734bf3308796ee18b", Balance: 1,
			PowNonce: 1, Address: "aaaaaaaaaa", Timestamp: 1, PubKey: "bbbbbbbb"}
		open.SetHash()

		tx, err := bs.Store(open)

		tx = &ledger.Transaction{Type: ledger.Transaction_SEND, Link: "dddddddddddddd", Previous: open.Hash,
			Signature: "df0d25f706c31d2007ed91da185ac727e5e38bc77f4309bb587e1ff7557ace39", Balance: 1,
			PowNonce: 1, Address: "aaaaaaaaaa", Timestamp: 1, PubKey: "bbbbbbbb",
		}
		tx.SetHash()

		b :=  &ledger.Transaction{}
		ms.Put("dddddddddddddd", b.ToBytes())

		tx, err = bs.Store(tx)

		Expect(err).To(BeNil())
		Expect(tx).NotTo(BeNil())

		chain, err := bs.GetTransactionChain(string(tx.Hash), true)
		Expect(err).To(BeNil())
		Expect(len(chain)).To(Equal(2))

		Expect(chain[0].Hash).To(Equal("e54ac468ad78c4b29262cd70331e7f1a53a67f28919419e2ee4134ea895a85bc"))
		Expect(chain[1].Hash).To(Equal("a04714dbf2c4ec17b38e6d1a2e15c70fcb90855e896b4e12c087c3e0af17c136"))
	})
})
