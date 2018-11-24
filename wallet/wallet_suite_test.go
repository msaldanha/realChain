package wallet_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestWallet(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Wallet Suite")
}
