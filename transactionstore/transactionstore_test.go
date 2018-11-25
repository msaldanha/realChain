package transactionstore_test

import (
	"github.com/msaldanha/realChain/protocol"
	"github.com/msaldanha/realChain/tests"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/msaldanha/realChain/transaction"
	"github.com/msaldanha/realChain/transactionstore"
	"github.com/golang/mock/gomock"
	"time"
)

var _ = Describe("TransactionStore", func() {

	It("Should not accept empty/partially filled transaction", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := tests.CreateNonEmptyMemoryStore()
		val := transaction.NewValidatorCreator()
		bs := transactionstore.New(ms, val)
		tx := &transaction.Transaction{}
		tx, err := bs.Store(tx)
		Expect(tx).To(BeNil())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(transaction.ErrInvalidTransactionType))

		tx = &transaction.Transaction{Type: protocol.Transaction_SEND, Link: []byte("ddddddddddddd"), Previous: []byte("ppppppppp"), Signature: []byte("ssssssss"), Balance: 1,
			PowNonce: 1, Address: []byte("aaaaaaaaaa"), Representative: []byte("rrrrrrrrrrrrrrr"), Timestamp: time.Now().Unix(), PubKey: []byte("kkkkkkk")}
		tx1, err := bs.Store(tx)
		Expect(tx1).To(BeNil())
		Expect(err).NotTo(BeNil())
		Expect(err).To(Equal(transaction.ErrTransactionHashCantBeEmpty))

		tx.SetHash()

		dest := &transaction.Transaction{Type: protocol.Transaction_OPEN, Link: []byte("ddddddddddddd"), Previous: []byte("ppppppppp"), Signature: []byte("ssssssss"), Balance: 1,
			PowNonce: 1, Address: []byte("aaaaaaaaaa"), Representative: []byte("rrrrrrrrrrrrrrr"), Timestamp: time.Now().Unix()}
		ms.Put("ddddddddddddd", dest.ToBytes())

		tx, err = bs.Store(tx)
		Expect(tx).NotTo(BeNil())
		Expect(err).To(BeNil())
	})

	It("Should accept properly filled transaction", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := tests.CreateNonEmptyMemoryStore()
		val := transaction.NewValidatorCreator()
		bs := transactionstore.New(ms, val)

		tx := &transaction.Transaction{Type: protocol.Transaction_SEND, Link: []byte("ddddddddddddd"), Previous: []byte("ppppppppp"),
			Signature: []byte("a246ce6b1d2b57ac33073127d8f9539fca32fb48481d46d734bf3308796ee18b"), Balance: 1,
			PowNonce: 1, Address: []byte("aaaaaaaaaa"), Representative: []byte("rrrrrrrrrrrrrrr"), Timestamp: 1, PubKey: []byte("kkkkkkk")}
		tx.SetHash()

		dest := &transaction.Transaction{Type: protocol.Transaction_OPEN, Link: []byte("ddddddddddddd"), Previous: []byte("ppppppppp"), Signature: []byte("ssssssss"), Balance: 1,
			PowNonce: 1, Address: []byte("aaaaaaaaaa"), Representative: []byte("rrrrrrrrrrrrrrr"), Timestamp: 1}
		ms.Put("ddddddddddddd", dest.ToBytes())

		tx, err := bs.Store(tx)
		Expect(err).To(BeNil())
		Expect(tx).NotTo(BeNil())

		txFromKeyStore, _, _ := ms.Get(string(tx.Hash))
		Expect(transaction.NewTransactionFromBytes(txFromKeyStore)).To(Equal(tx))
	})

	It("Should extract the chain for the transaction", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		ms := tests.CreateNonEmptyMemoryStore()
		val := transaction.NewValidatorCreator()
		bs := transactionstore.New(ms, val)

		open := &transaction.Transaction{Type: protocol.Transaction_OPEN, Link: []byte("ddddddddddddd"), Previous: []byte("ppppppppp"),
			Signature: []byte("a246ce6b1d2b57ac33073127d8f9539fca32fb48481d46d734bf3308796ee18b"), Balance: 1,
			PowNonce: 1, Address: []byte("aaaaaaaaaa"), Representative: []byte("rrrrrrrrrrrrrrr"), Timestamp: 1, PubKey: []byte("kkkkkkk")}
		open.SetHash()

		tx, err := bs.Store(open)

		tx = &transaction.Transaction{Type: protocol.Transaction_SEND, Link: []byte("ddddddddddddd"), Previous: open.Hash,
			Signature: []byte("df0d25f706c31d2007ed91da185ac727e5e38bc77f4309bb587e1ff7557ace39"), Balance: 1,
			PowNonce: 1, Address: []byte("aaaaaaaaaa"), Representative: []byte("rrrrrrrrrrrrrrr"), Timestamp: 1, PubKey: []byte("kkkkkkk"),
		}
		tx.SetHash()

		b :=  &transaction.Transaction{}
		ms.Put("ddddddddddddd", b.ToBytes())

		tx, err = bs.Store(tx)

		Expect(err).To(BeNil())
		Expect(tx).NotTo(BeNil())

		chain, err := bs.GetTransactionChain(string(tx.Hash), true)
		Expect(err).To(BeNil())
		Expect(len(chain)).To(Equal(2))

		Expect(chain[0].Hash).To(Equal([]byte("a246ce6b1d2b57ac33073127d8f9539fca32fb48481d46d734bf3308796ee18b")))
		Expect(chain[1].Hash).To(Equal([]byte("df0d25f706c31d2007ed91da185ac727e5e38bc77f4309bb587e1ff7557ace39")))
	})
})
