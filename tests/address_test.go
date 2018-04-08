package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/golang/mock/gomock"
	"github.com/msaldanha/realChain/address"
	"fmt"
	"encoding/hex"
)

var _ = Describe("Address", func() {
	It("Should create an address", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		//privateKey := "e3d6fd07b55459dbde1c9439fd2075fb76c9c0427562e6ae17aefde74035d40d"
		pubKey := "02edffb87f10eb61ee4dbd1ac3a7c80477ca515682d7651b6d59e2fc29f20290d50e6d0fb705b74c884d1961e19babac2580d895ee87eec4c3fa49a692d5a027"
		expectedAddr := "1PEY9rskiiiX4tPUXHjZYuV9qepriaxgqJ"

		addr := address.New()
		dec, _ := hex.DecodeString(pubKey)
		ad, err := addr.GenerateForKey(dec)
		Expect(err).To(BeNil())
		Expect(ad).To(Equal(expectedAddr))
		fmt.Println(ad)
	})

	It("Should validate an address", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		expectedAddr := "1PEY9rskiiiX4tPUXHjZYuV9qepriaxgqJ"

		addr := address.New()
		ok, err := addr.IsValid(expectedAddr)
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())

		expectedAddr = "2PEY9rskiiiX4tPUXHjZYuV9qepriaxgqJ"

		ok, err = addr.IsValid(expectedAddr)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("invalid checksum"))
		Expect(ok).To(BeFalse())
	})
})
