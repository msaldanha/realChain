package block

import (
	"fmt"
	"strconv"
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/gob"
)

type Type int16

const (
	OPEN    Type = 1 + iota
	SEND
	RECEIVE
	CHANGE
)

type Block struct {
	Timestamp      int64
	Type           Type
	Account        []byte
	Representative []byte
	Previous       []byte
	Link           []byte
	Balance        float64
	Hash           []byte
	Signature      []byte
	PowTarget      int16
	PowNonce       int64
	PubKey         []byte
}

func (bt Type) IsValid() (bool) {
	return bt >= OPEN && bt <= CHANGE
}

func (bt Type) String() (string) {
	name := ""
	switch bt {
	case OPEN:
		name = "OPEN"
	case SEND:
		name = "SEND"
	case RECEIVE:
		name = "RECEIVE"
	case CHANGE:
		name = "CHANGE"
	}
	return fmt.Sprintf("%d(%s)", int(bt), name)
}

func (b *Block) SetHash() (error) {
	hash, err := b.GetHash()
	if err != nil {
		return err
	}
	b.Hash = hash
	return nil
}

func (b *Block) GetHash() ([]byte, error) {
	hashableBytes, err := b.GetHashableBytes()
	if err != nil {
		return nil, err
	}
	headers := bytes.Join(hashableBytes, []byte{})
	hash := sha256.Sum256(headers)
	return []byte(hex.EncodeToString(hash[:])), nil
}

func (b *Block) GetHashableBytes() ([][]byte, error) {
	var balance bytes.Buffer
	if err := binary.Write(&balance, binary.LittleEndian, b.Balance); err != nil {
		return nil, err
	}
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	return [][]byte{timestamp, b.Account, b.Representative,
		b.Previous, b.Link, balance.Bytes()}, nil
}

func (b *Block) ToBytes() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	encoder.Encode(b)
	return result.Bytes()
}

func NewBlockFromBytes(d []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(d))
	decoder.Decode(&block)
	return &block
}