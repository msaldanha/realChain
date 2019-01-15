package ledger

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"github.com/golang/protobuf/proto"
	"github.com/msaldanha/realChain/crypto"
	"log"
	"math"
	"math/big"
	"strconv"
	"time"
)

//go:generate protoc -I.. ledger/transaction.proto --go_out=plugins=grpc:../

const targetBits int16 = 16

func NewOpenTransaction() *Transaction {
	return &Transaction{Type: Transaction_OPEN, Timestamp: time.Now().UnixNano()}
}

func NewSendTransaction() *Transaction {
	return &Transaction{Type: Transaction_SEND, Timestamp: time.Now().UnixNano()}
}

func NewReceiveTransaction() *Transaction {
	return &Transaction{Type: Transaction_RECEIVE, Timestamp: time.Now().UnixNano()}
}

func (tx *Transaction) SetHash() error {
	hash, err := tx.CalculateHash()
	if err != nil {
		return err
	}
	tx.Hash = hash
	return nil
}

func (tx *Transaction) CalculateHash() ([]byte, error) {
	hashableBytes, err := tx.GetHashableBytes()
	if err != nil {
		return nil, err
	}
	headers := bytes.Join(hashableBytes, []byte{})
	hash := sha256.Sum256(headers)
	return []byte(hex.EncodeToString(hash[:])), nil
}

func (tx *Transaction) GetHashableBytes() ([][]byte, error) {
	var balance bytes.Buffer
	if err := binary.Write(&balance, binary.LittleEndian, tx.Balance); err != nil {
		return nil, err
	}
	timestamp := []byte(strconv.FormatInt(tx.Timestamp, 10))
	return [][]byte{timestamp, tx.Address, tx.Representative,
		tx.Previous, tx.Link, balance.Bytes()}, nil
}

func (tx *Transaction) CalculatePow() (int64, []byte, error) {
	var hashInt big.Int
	var hash [32]byte
	var nonce int64 = 0

	target := getTarget()

	data, err := tx.GetHashableBytes()
	if err != nil {
		return 0, nil, err
	}

	for nonce < math.MaxInt64 {
		dataWithNonce := append(data, int64ToBytes(nonce))
		hash = sha256.Sum256(bytes.Join(dataWithNonce, []byte{}))
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(target) == -1 {
			break
		} else {
			nonce++
		}
	}

	hexHash := []byte(hex.EncodeToString(hash[:]))

	return nonce, hexHash[:], nil
}

func (tx *Transaction) SetPow() error {
	nonce, hash, err := tx.CalculatePow()
	if err != nil {
		return err
	}
	tx.PowNonce = nonce
	tx.Hash = hash
	return nil
}

func (tx *Transaction) VerifyPow() (bool, error) {
	var hashInt big.Int

	target := getTarget()

	data, err := tx.GetHashableBytes()
	if err != nil {
		return false, err
	}
	dataWithNonce := append(data, int64ToBytes(tx.PowNonce))
	hash := sha256.Sum256(bytes.Join(dataWithNonce, []byte{}))
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(target) == -1, nil
}

func (tx *Transaction) Sign(privateKey *ecdsa.PrivateKey) error {
	s, err := crypto.Sign(tx.Hash, privateKey)
	if err != nil {
		return err
	}
	tx.Signature = s
	return nil
}

func (tx *Transaction) VerifySignature() bool {
	return crypto.VerifySignature(tx.Signature, tx.PubKey, tx.Hash)
}

func (tx *Transaction) ToBytes() []byte {
	data, _ := proto.Marshal(tx)
	return data
}

func NewTransactionFromBytes(d []byte) *Transaction {
	tx := &Transaction{}
	_ = proto.Unmarshal(d, tx)
	return tx
}

func getTarget() *big.Int {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
	return target
}

func int64ToBytes(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
