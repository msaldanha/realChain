package tests

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
			PowNonce: 1, Account: []byte("aaaaaaaaaa"), Representative: []byte("rrrrrrrrrrrrrrr"),
			Timestamp: 1}

		tx.SetHash()
		Expect(tx.Hash).To(Equal([]byte("a246ce6b1d2b57ac33073127d8f9539fca32fb48481d46d734bf3308796ee18b")))
	})
})
