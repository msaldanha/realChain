package block

import (
	"fmt"
	"strconv"
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
)

type BlockType int16

const (
	OPEN BlockType = 1 + iota
	SEND
	RECEIVE
	CHANGE
)

type Block struct {
	Timestamp     	int64
	Type 			BlockType
	Account 		[]byte
	Representative 	[]byte
	Previous 		[]byte
	Link 			[]byte
	Work 			[]byte
	Balance 		float64
	Hash 			[]byte
	Signature 		[]byte
}

func (bt BlockType) IsValid() (bool) {
	return bt >= OPEN && bt <= CHANGE
}

func (bt BlockType) String() (string) {
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
	var balance bytes.Buffer
	if err := binary.Write(&balance, binary.LittleEndian, b.Balance); err != nil {
		return err
	}
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{timestamp, b.Account, b.Representative,
		b.Previous, b.Link, b.Work, balance.Bytes()}, []byte{})
	hash := sha256.Sum256(headers)
	b.Hash = []byte(hex.EncodeToString(hash[:]))
	return nil
}