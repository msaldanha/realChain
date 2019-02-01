package ledger_test

import (
	"github.com/msaldanha/realChain/ledger"
	"github.com/msaldanha/realChain/tests"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/golang/mock/gomock"
	"time"
)

var _ = Describe("Validator", func() {
	It("Should not accept empty/partial filled transaction for OPEN type", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := tests.CreateNonEmptyMemoryStore()
		val := ledger.NewValidatorCreator().CreateValidatorForTransaction(ledger.Transaction_OPEN, ms)

		tx := &ledger.Transaction{}
		ok, err := val.IsFilled(tx)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrInvalidTransactionType))

		tx.Type = ledger.Transaction_OPEN

		ok, err = val.IsFilled(tx)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrInvalidTransactionTimestamp))

		tx.Timestamp = time.Now().Unix()

		ok, err = val.IsFilled(tx)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrTransactionAddressCantBeEmpty))

		tx.Address = "dddddddddddd"

		ok, err = val.IsFilled(tx)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrTransactionSignatureCantBeEmpty))

		tx.Signature = "aaaaaaaaaa"

		ok, err = val.IsFilled(tx)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrTransactionPowNonceCantBeZero))

		tx.PowNonce = 1

		ok, err = val.IsFilled(tx)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrTransactionHashCantBeEmpty))

		tx.Hash = "dddddddddddd"

		ok, err = val.IsFilled(tx)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrPubKeyCantBeEmpty))

		tx.PubKey = "dddddddddddd"

		ok, err = val.IsFilled(tx)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrTransactionLinkCantBeEmpty))

		tx.Link = "dddddddddddd"

		ok, err = val.IsFilled(tx)
		Expect(ok).To(BeTrue())
		Expect(err).To(BeNil())
	})

	It("Should not accept empty/partial filled transaction for SEND type", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := tests.CreateNonEmptyMemoryStore()
		val := ledger.NewValidatorCreator().CreateValidatorForTransaction(ledger.Transaction_SEND, ms)

		tx := &ledger.Transaction{}
		ok, err := val.IsFilled(tx)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrInvalidTransactionType))

		tx.Type = ledger.Transaction_SEND

		tests.AssertCommonVal(val, tx)

		ok, err = val.IsFilled(tx)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrTransactionLinkCantBeEmpty))

		tx.Link = "dddddddddddd"
	})

	It("Should not accept empty/partial filled transaction for RECEIVE type", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := tests.CreateNonEmptyMemoryStore()
		val := ledger.NewValidatorCreator().CreateValidatorForTransaction(ledger.Transaction_RECEIVE, ms)

		tx := &ledger.Transaction{}
		ok, err := val.IsFilled(tx)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrInvalidTransactionType))

		tx.Type = ledger.Transaction_RECEIVE

		tests.AssertCommonVal(val, tx)

		ok, err = val.IsFilled(tx)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrTransactionLinkCantBeEmpty))

		tx.Link = "dddddddddddd"
	})

	It("Should not accept OPEN transaction with invalid fields", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := tests.CreateNonEmptyMemoryStore()
		val := ledger.NewValidatorCreator().CreateValidatorForTransaction(ledger.Transaction_OPEN, ms)
		tx := &ledger.Transaction{Type: ledger.Transaction_OPEN, Link: "dddddddddddd", Previous: "bbbbbbbbbb",
		Signature: "ffffffffff", Balance: 1, Timestamp: 1, PubKey: "eeeeeeeeee",
			PowNonce: 1}
		tx.SetHash()

		ok, err := val.IsValid(tx)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrTransactionAddressCantBeEmpty))

		tx.Address = "aaaaaaaaaa"

		ok, err = val.IsValid(tx)
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())
	})

	It("Should not accept SEND transaction with invalid fields", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := tests.CreateNonEmptyMemoryStore()
		val := ledger.NewValidatorCreator().CreateValidatorForTransaction(ledger.Transaction_SEND, ms)
		tx := &ledger.Transaction{Type: ledger.Transaction_SEND, Previous: "bbbbbbbbbb", Signature: "ffffffffff", Balance: 1,
			PowNonce: 1, Address: "aaaaaaaaaa", Timestamp: 1, PubKey: "eeeeeeeeee"}
		tx.SetHash()

		ok, err := val.IsValid(tx)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrTransactionLinkCantBeEmpty))

		tx.Link = "dddddddddddd"

		ok, err = val.IsValid(tx)
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())

	})

	It("Should not accept RECEIVE transaction with invalid fields", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := tests.CreateNonEmptyMemoryStore()
		val := ledger.NewValidatorCreator().CreateValidatorForTransaction(ledger.Transaction_RECEIVE, ms)
		tx := &ledger.Transaction{Type: ledger.Transaction_RECEIVE, Link: "dddddddddddd", Previous: "bbbbbbbbbb", Signature: "ffffffffff", Balance: 1,
			PowNonce: 1, Address: "aaaaaaaaaa", Timestamp: 1, PubKey: "eeeeeeeeee"}
		tx.SetHash()

		ok, err := val.IsValid(tx)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrSourceNotFound))

		source := &ledger.Transaction{Type: ledger.Transaction_OPEN, Link: "dddddddddddd", Previous: "bbbbbbbbbb", Signature: "ffffffffff", Balance: 1,
			PowNonce: 1, Address: "aaaaaaaaaa", Timestamp: 1, PubKey: "eeeeeeeeee"}
		ms.Put("dddddddddddd", source.ToBytes())

		ok, err = val.IsValid(tx)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrInvalidSourceType))

		source.Type = ledger.Transaction_SEND
		ms.Put("dddddddddddd", source.ToBytes())

		ok, err = val.IsValid(tx)
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())
	})

	It("Should not accept CHANGE transaction with invalid fields", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := tests.CreateNonEmptyMemoryStore()
		val := ledger.NewValidatorCreator().CreateValidatorForTransaction(ledger.Transaction_CHANGE, ms)
		tx := &ledger.Transaction{Type: ledger.Transaction_CHANGE, Link: "dddddddddddd", Previous: "bbbbbbbbbb", Signature: "ffffffffff", Balance: 1,
			PowNonce: 1, Address: "aaaaaaaaaa", Timestamp: 1, PubKey: "eeeeeeeeee"}
		tx.SetHash()

		ok, err := val.IsValid(tx)
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())
	})

})
