package ledger

import (
	"github.com/msaldanha/realChain/keypair"
	"bytes"
	"encoding/gob"
)

type Account struct {
	Keys *keypair.KeyPair
	Address string
}

func (acc *Account) ToBytes() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	encoder.Encode(acc)
	return result.Bytes()
}

func NewAccountFromBytes(a []byte) *Account {
	var acc Account
	decoder := gob.NewDecoder(bytes.NewReader(a))
	decoder.Decode(&acc)
	return &acc
}