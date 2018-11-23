package keyvaluestore_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestKeyvaluestore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Keyvaluestore Suite")
}
