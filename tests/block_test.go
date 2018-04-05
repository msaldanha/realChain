package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/golang/mock/gomock"
	. "github.com/msaldanha/realChain/block"
)

var _ = Describe("Block", func() {
	It("Should calculate its hash", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		block := &Block{Type: RECEIVE, Link: []byte("ddddddddddddd"), Previous: []byte("ppppppppp"),
			Signature: []byte("ssssssss"), Balance: 1,
			Work: []byte("wwwwwwww"), Account: []byte("aaaaaaaaaa"), Representative: []byte("rrrrrrrrrrrrrrr"),
			Timestamp: 1}

		block.SetHash()
		Expect(block.Hash).To(Equal([]byte("bcf58e48dabd2937f86396d671910425d1b4acf3a6ff9ac95c3babb84ca9f0ad")))
	})
})
