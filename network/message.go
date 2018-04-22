package network

import (
	"bytes"
	"github.com/davecgh/go-xdr/xdr2"
)

const (
	Magic int32 = 0xcaba
	Version int32 = 0x0
)

type message struct {
	Magic    int32
	Version  int32
	EndPoint string
	Payload  []byte
}

func NewMessage() *message {
	return &message{Magic:Magic, Version:Version}
}

func (b *message) ToBytes() []byte {
	var result bytes.Buffer
	encoder := xdr.NewEncoder(&result)
	encoder.Encode(b)
	return result.Bytes()
}
