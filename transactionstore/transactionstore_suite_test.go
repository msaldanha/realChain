package transactionstore_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTransactionstore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Transactionstore Suite")
}
