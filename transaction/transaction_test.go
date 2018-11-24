package transaction_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/golang/mock/gomock"
	. "github.com/msaldanha/realChain/transaction"
)

var _ = Describe("Transaction", func() {
	It("Should calculate its hash", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		tx := &Transaction{Type: RECEIVE, Link: []byte("ddddddddddddd"), Previous: []byte("ppppppppp"),
			Signature: []byte("ssssssss"), Balance: 1,
			PowNonce: 1, Address: []byte("aaaaaaaaaa"), Representative: []byte("rrrrrrrrrrrrrrr"),
			Timestamp: 1}

		tx.SetHash()
		Expect(tx.Hash).To(Equal([]byte("a246ce6b1d2b57ac33073127d8f9539fca32fb48481d46d734bf3308796ee18b")))
	})

	It("Should calculate the PoW ", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		tx := &Transaction{Type: SEND, Link: []byte("ddddddddddddd"), Previous: []byte("ppppppppp"),
			Signature: []byte("777d713768de05cb16cbc24eef83b43b20a3a80dce05549f130aaf5a4234e4c2"), Balance: 1,
			PowNonce: 1, Address: []byte("aaaaaaaaaa"), Representative: []byte("rrrrrrrrrrrrrrr"), Timestamp: 1}
		tx.SetHash()

		err := tx.SetPow()
		Expect(err).To(BeNil())

		Expect(tx.PowNonce).To(Equal(int64(33794)))
		Expect(tx.Hash).To(Equal([]byte("0000f4722f6416ddb43a4ee56921dd3a24c93b051a570e14ca07cd174517cf12")))
	})

	It("Should verify the PoW ", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		tx := &Transaction{Type: SEND, Link: []byte("ddddddddddddd"), Previous: []byte("ppppppppp"),
			Signature: []byte("777d713768de05cb16cbc24eef83b43b20a3a80dce05549f130aaf5a4234e4c2"), Balance: 1,
			PowNonce: 1, Address: []byte("aaaaaaaaaa"), Representative: []byte("rrrrrrrrrrrrrrr"), Timestamp: 1}
		tx.SetHash()

		tx.PowNonce = int64(33794)
		tx.Hash = []byte("0000f4722f6416ddb43a4ee56921dd3a24c93b051a570e14ca07cd174517cf12")

		ok, err := tx.VerifyPow()
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())

	})
})
