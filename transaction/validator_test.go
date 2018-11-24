package transaction_test

import (
	"github.com/msaldanha/realChain/tests"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/golang/mock/gomock"
	"github.com/msaldanha/realChain/transaction"
	"time"
)

var _ = Describe("Validator", func() {
	It("Should not accept empty/partial filled transaction for OPEN type", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := tests.CreateNonEmptyMemoryStore()
		val := transaction.NewValidatorCreator().CreateValidatorForTransaction(transaction.OPEN, ms)

		tx := &transaction.Transaction{}
		ok, err := val.IsFilled(tx)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(transaction.ErrInvalidTransactionType))

		tx.Type = transaction.OPEN

		ok, err = val.IsFilled(tx)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(transaction.ErrInvalidTransactionTimestamp))

		tx.Timestamp = time.Now().Unix()

		ok, err = val.IsFilled(tx)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(transaction.ErrTransactionAddressCantBeEmpty))

		tx.Address = []byte("xxxxxxxxxxxxxxxxxxx")

		ok, err = val.IsFilled(tx)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(transaction.ErrTransactionSignatureCantBeEmpty))

		tx.Signature = []byte("ssssssssssssssssssssss")

		ok, err = val.IsFilled(tx)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(transaction.ErrTransactionPowNonceCantBeZero))

		tx.PowNonce = 1

		ok, err = val.IsFilled(tx)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(transaction.ErrTransactionHashCantBeEmpty))

		tx.Hash = []byte("xxxxxxxxxxxxxxxxxxx")

		ok, err = val.IsFilled(tx)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(transaction.ErrPubKeyCantBeEmpty))

		tx.PubKey = []byte("xxxxxxxxxxxxxxxxxxx")

		ok, err = val.IsFilled(tx)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(transaction.ErrTransactionLinkCantBeEmpty))

		tx.Link = []byte("xxxxxxxxxxxxxxxxxxx")

		ok, err = val.IsFilled(tx)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(transaction.ErrTransactionRepresentativeCantBeEmpty))

		tx.Representative = []byte("xxxxxxxxxxxxxxxxxxx")

		ok, err = val.IsFilled(tx)
		Expect(ok).To(BeTrue())
		Expect(err).To(BeNil())
	})

	It("Should not accept empty/partial filled transaction for SEND type", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := tests.CreateNonEmptyMemoryStore()
		val := transaction.NewValidatorCreator().CreateValidatorForTransaction(transaction.SEND, ms)

		tx := &transaction.Transaction{}
		ok, err := val.IsFilled(tx)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(transaction.ErrInvalidTransactionType))

		tx.Type = transaction.SEND

		tests.AssertCommonVal(val, tx)

		ok, err = val.IsFilled(tx)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(transaction.ErrTransactionLinkCantBeEmpty))

		tx.Link = []byte("xxxxxxxxxxxxxxxxxxx")
	})

	It("Should not accept empty/partial filled transaction for RECEIVE type", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := tests.CreateNonEmptyMemoryStore()
		val := transaction.NewValidatorCreator().CreateValidatorForTransaction(transaction.RECEIVE, ms)

		tx := &transaction.Transaction{}
		ok, err := val.IsFilled(tx)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(transaction.ErrInvalidTransactionType))

		tx.Type = transaction.RECEIVE

		tests.AssertCommonVal(val, tx)

		ok, err = val.IsFilled(tx)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(transaction.ErrTransactionLinkCantBeEmpty))

		tx.Link = []byte("xxxxxxxxxxxxxxxxxxx")
	})

	It("Should not accept empty/partial filled transaction for CHANGE type", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := tests.CreateNonEmptyMemoryStore()
		val := transaction.NewValidatorCreator().CreateValidatorForTransaction(transaction.CHANGE, ms)

		tx := &transaction.Transaction{}
		ok, err := val.IsFilled(tx)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(transaction.ErrInvalidTransactionType))

		tx.Type = transaction.CHANGE

		tests.AssertCommonVal(val, tx)

		ok, err = val.IsFilled(tx)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(transaction.ErrTransactionRepresentativeCantBeEmpty))

		tx.Representative = []byte("xxxxxxxxxxxxxxxxxxx")
	})

	It("Should not accept OPEN transaction with invalid fields", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := tests.CreateNonEmptyMemoryStore()
		val := transaction.NewValidatorCreator().CreateValidatorForTransaction(transaction.OPEN, ms)
		tx := &transaction.Transaction{Type: transaction.OPEN, Link: []byte([]byte("ddddddddddddd")), Previous: []byte([]byte("ppppppppp")),
		Signature: []byte([]byte("ssssssss")), Balance: 1, Timestamp: 1, PubKey: []byte("kkkkkkk"),
			PowNonce: 1, Representative:[]byte("rrrrrrrrrrrrrrr")}
		tx.SetHash()

		ok, err := val.IsValid(tx)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(transaction.ErrTransactionAddressCantBeEmpty))

		tx.Address = []byte("aaaaaaaaaa")

		ok, err = val.IsValid(tx)
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())
	})

	It("Should not accept SEND transaction with invalid fields", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := tests.CreateNonEmptyMemoryStore()
		val := transaction.NewValidatorCreator().CreateValidatorForTransaction(transaction.SEND, ms)
		tx := &transaction.Transaction{Type: transaction.SEND, Previous: []byte("ppppppppp"), Signature: []byte("ssssssss"), Balance: 1,
			PowNonce: 1, Address:[]byte("aaaaaaaaaa"), Representative:[]byte("rrrrrrrrrrrrrrr"), Timestamp: 1, PubKey: []byte("kkkkkkk")}
		tx.SetHash()

		ok, err := val.IsValid(tx)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(transaction.ErrTransactionLinkCantBeEmpty))

		tx.Link = []byte("ddddddddddddd")

		ok, err = val.IsValid(tx)
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())

	})

	It("Should not accept RECEIVE transaction with invalid fields", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := tests.CreateNonEmptyMemoryStore()
		val := transaction.NewValidatorCreator().CreateValidatorForTransaction(transaction.RECEIVE, ms)
		tx := &transaction.Transaction{Type: transaction.RECEIVE, Link: []byte("ddddddddddddd"), Previous: []byte("ppppppppp"), Signature: []byte("ssssssss"), Balance: 1,
			PowNonce: 1, Address:[]byte("aaaaaaaaaa"), Representative:[]byte("rrrrrrrrrrrrrrr"), Timestamp: 1, PubKey: []byte("kkkkkkk")}
		tx.SetHash()

		ok, err := val.IsValid(tx)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(transaction.ErrSourceNotFound))

		source := &transaction.Transaction{Type: transaction.OPEN, Link: []byte("ddddddddddddd"), Previous: []byte("ppppppppp"), Signature: []byte("ssssssss"), Balance: 1,
			PowNonce: 1, Address:[]byte("aaaaaaaaaa"), Representative:[]byte("rrrrrrrrrrrrrrr"), Timestamp: 1, PubKey: []byte("kkkkkkk")}
		ms.Put("ddddddddddddd", source.ToBytes())

		ok, err = val.IsValid(tx)
		Expect(ok).To(BeFalse())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(transaction.ErrInvalidSourceType))

		source.Type = transaction.SEND
		ms.Put("ddddddddddddd", source.ToBytes())

		ok, err = val.IsValid(tx)
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())
	})

	It("Should not accept CHANGE transaction with invalid fields", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := tests.CreateNonEmptyMemoryStore()
		val := transaction.NewValidatorCreator().CreateValidatorForTransaction(transaction.CHANGE, ms)
		tx := &transaction.Transaction{Type: transaction.CHANGE, Link: []byte("ddddddddddddd"), Previous: []byte("ppppppppp"), Signature: []byte("ssssssss"), Balance: 1,
			PowNonce: 1, Address:[]byte("aaaaaaaaaa"), Representative:[]byte("rrrrrrrrrrrrrrr"), Timestamp: 1, PubKey: []byte("kkkkkkk")}
		tx.SetHash()

		ok, err := val.IsValid(tx)
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())
	})

})
