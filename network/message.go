package network

import (
	"bytes"
	"github.com/davecgh/go-xdr/xdr2"
)

const (
	Magic int32 = 0xcaba
	Version int32 = 0x0
)

type Message struct {
	Magic    int32
	Version  int32
	EndPoint string
	Payload  []byte
}

func NewMessage() *Message {
	return &Message{Magic:Magic, Version:Version}
}


func NewMessageFromBytes(d []byte) *Message {
	var msg Message
	decoder := xdr.NewDecoder(bytes.NewReader(d))
	decoder.Decode(&msg)
	return &msg
}

func (b *Message) ToBytes() []byte {
	var result bytes.Buffer
	encoder := xdr.NewEncoder(&result)
	encoder.Encode(b)
	return result.Bytes()
}
