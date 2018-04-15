package ledger

import (
	"github.com/msaldanha/realChain/keypair"
	"bytes"
	"encoding/gob"
)

type Address struct {
	Keys *keypair.KeyPair
	Address string
}

func (a *Address) ToBytes() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	encoder.Encode(a)
	return result.Bytes()
}

func NewAddressFromBytes(a []byte) *Address {
	var acc Address
	decoder := gob.NewDecoder(bytes.NewReader(a))
	decoder.Decode(&acc)
	return &acc
}