package ledger_test

import (
	"github.com/golang/mock/gomock"
	"github.com/msaldanha/realChain/ledger"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Transaction", func() {
	It("Should calculate its hash", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		tx := &ledger.Transaction{Type: ledger.Transaction_RECEIVE, Link: "dddddddddddd", Previous: "bbbbbbbbbb",
			Signature: "ffffffffff", Balance: 1,
			PowNonce: 1, Address: "aaaaaaaaaa",
			Timestamp: 1}

		tx.SetHash()
		Expect(tx.Hash).To(Equal("86ae3bd95c4f1a200e3a8ee1655710321b2a2becaccd9ce6ed832cd4dc92502f"))
	})

	It("Should calculate the PoW ", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		tx := &ledger.Transaction{Type: ledger.Transaction_SEND, Link: "dddddddddddd", Previous: "bbbbbbbbbb",
			Signature: "777d713768de05cb16cbc24eef83b43b20a3a80dce05549f130aaf5a4234e4c2", Balance: 1,
			PowNonce: 1, Address: "aaaaaaaaaa", Timestamp: 1}
		tx.SetHash()

		err := tx.SetPow()
		Expect(err).To(BeNil())

		Expect(tx.PowNonce).To(Equal(int64(54914)))
		Expect(tx.Hash).To(Equal("0000a52076206cee2400d1a05a60d1768183ba8ddfd9ba6140ff13296f53dbaa"))
	})

	It("Should verify the PoW ", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		tx := &ledger.Transaction{Type: ledger.Transaction_SEND, Link: "dddddddddddd", Previous: "bbbbbbbbbb",
			Signature: "777d713768de05cb16cbc24eef83b43b20a3a80dce05549f130aaf5a4234e4c2", Balance: 1,
			PowNonce: 1, Address: "aaaaaaaaaa", Timestamp: 1}
		tx.SetHash()

		tx.PowNonce = int64(54914)
		tx.Hash = "00006a216f8ea0ecc5f106bd6401b7986a0e45670fc55aa2fbe6b49c3dbb4133"

		ok, err := tx.VerifyPow()
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())

	})
})
