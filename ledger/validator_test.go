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

		tx.Address = []byte("xxxxxxxxxxxxxxxxxxx")

		ok, err = val.IsFilled(tx)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrTransactionSignatureCantBeEmpty))

		tx.Signature = []byte("ssssssssssssssssssssss")

		ok, err = val.IsFilled(tx)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrTransactionPowNonceCantBeZero))

		tx.PowNonce = 1

		ok, err = val.IsFilled(tx)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrTransactionHashCantBeEmpty))

		tx.Hash = []byte("xxxxxxxxxxxxxxxxxxx")

		ok, err = val.IsFilled(tx)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrPubKeyCantBeEmpty))

		tx.PubKey = []byte("xxxxxxxxxxxxxxxxxxx")

		ok, err = val.IsFilled(tx)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrTransactionLinkCantBeEmpty))

		tx.Link = []byte("xxxxxxxxxxxxxxxxxxx")

		ok, err = val.IsFilled(tx)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrTransactionRepresentativeCantBeEmpty))

		tx.Representative = []byte("xxxxxxxxxxxxxxxxxxx")

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

		tx.Link = []byte("xxxxxxxxxxxxxxxxxxx")
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

		tx.Link = []byte("xxxxxxxxxxxxxxxxxxx")
	})

	It("Should not accept empty/partial filled transaction for CHANGE type", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := tests.CreateNonEmptyMemoryStore()
		val := ledger.NewValidatorCreator().CreateValidatorForTransaction(ledger.Transaction_CHANGE, ms)

		tx := &ledger.Transaction{}
		ok, err := val.IsFilled(tx)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrInvalidTransactionType))

		tx.Type = ledger.Transaction_CHANGE

		tests.AssertCommonVal(val, tx)

		ok, err = val.IsFilled(tx)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrTransactionRepresentativeCantBeEmpty))

		tx.Representative = []byte("xxxxxxxxxxxxxxxxxxx")
	})

	It("Should not accept OPEN transaction with invalid fields", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := tests.CreateNonEmptyMemoryStore()
		val := ledger.NewValidatorCreator().CreateValidatorForTransaction(ledger.Transaction_OPEN, ms)
		tx := &ledger.Transaction{Type: ledger.Transaction_OPEN, Link: []byte([]byte("ddddddddddddd")), Previous: []byte([]byte("ppppppppp")),
		Signature: []byte([]byte("ssssssss")), Balance: 1, Timestamp: 1, PubKey: []byte("kkkkkkk"),
			PowNonce: 1, Representative:[]byte("rrrrrrrrrrrrrrr")}
		tx.SetHash()

		ok, err := val.IsValid(tx)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrTransactionAddressCantBeEmpty))

		tx.Address = []byte("aaaaaaaaaa")

		ok, err = val.IsValid(tx)
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())
	})

	It("Should not accept SEND transaction with invalid fields", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := tests.CreateNonEmptyMemoryStore()
		val := ledger.NewValidatorCreator().CreateValidatorForTransaction(ledger.Transaction_SEND, ms)
		tx := &ledger.Transaction{Type: ledger.Transaction_SEND, Previous: []byte("ppppppppp"), Signature: []byte("ssssssss"), Balance: 1,
			PowNonce: 1, Address:[]byte("aaaaaaaaaa"), Representative:[]byte("rrrrrrrrrrrrrrr"), Timestamp: 1, PubKey: []byte("kkkkkkk")}
		tx.SetHash()

		ok, err := val.IsValid(tx)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrTransactionLinkCantBeEmpty))

		tx.Link = []byte("ddddddddddddd")

		ok, err = val.IsValid(tx)
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())

	})

	It("Should not accept RECEIVE transaction with invalid fields", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := tests.CreateNonEmptyMemoryStore()
		val := ledger.NewValidatorCreator().CreateValidatorForTransaction(ledger.Transaction_RECEIVE, ms)
		tx := &ledger.Transaction{Type: ledger.Transaction_RECEIVE, Link: []byte("ddddddddddddd"), Previous: []byte("ppppppppp"), Signature: []byte("ssssssss"), Balance: 1,
			PowNonce: 1, Address:[]byte("aaaaaaaaaa"), Representative:[]byte("rrrrrrrrrrrrrrr"), Timestamp: 1, PubKey: []byte("kkkkkkk")}
		tx.SetHash()

		ok, err := val.IsValid(tx)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrSourceNotFound))

		source := &ledger.Transaction{Type: ledger.Transaction_OPEN, Link: []byte("ddddddddddddd"), Previous: []byte("ppppppppp"), Signature: []byte("ssssssss"), Balance: 1,
			PowNonce: 1, Address:[]byte("aaaaaaaaaa"), Representative:[]byte("rrrrrrrrrrrrrrr"), Timestamp: 1, PubKey: []byte("kkkkkkk")}
		ms.Put("ddddddddddddd", source.ToBytes())

		ok, err = val.IsValid(tx)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(ledger.ErrInvalidSourceType))

		source.Type = ledger.Transaction_SEND
		ms.Put("ddddddddddddd", source.ToBytes())

		ok, err = val.IsValid(tx)
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())
	})

	It("Should not accept CHANGE transaction with invalid fields", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := tests.CreateNonEmptyMemoryStore()
		val := ledger.NewValidatorCreator().CreateValidatorForTransaction(ledger.Transaction_CHANGE, ms)
		tx := &ledger.Transaction{Type: ledger.Transaction_CHANGE, Link: []byte("ddddddddddddd"), Previous: []byte("ppppppppp"), Signature: []byte("ssssssss"), Balance: 1,
			PowNonce: 1, Address:[]byte("aaaaaaaaaa"), Representative:[]byte("rrrrrrrrrrrrrrr"), Timestamp: 1, PubKey: []byte("kkkkkkk")}
		tx.SetHash()

		ok, err := val.IsValid(tx)
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())
	})

})
