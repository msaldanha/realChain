package transaction

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"github.com/davecgh/go-xdr/xdr2"
	"github.com/msaldanha/realChain/protocol"
	"log"
	"math"
	"math/big"
	"strconv"
	"time"
)

const targetBits int16 = 16

type Transaction protocol.Transaction
type Type protocol.Transaction_Type

func NewOpenTransaction() *Transaction {
	return &Transaction{Type: protocol.Transaction_OPEN, Timestamp: time.Now().Unix()}
}

func NewSendTransaction() *Transaction {
	return &Transaction{Type: protocol.Transaction_SEND, Timestamp: time.Now().Unix()}
}

func NewReceiveTransaction() *Transaction {
	return &Transaction{Type: protocol.Transaction_RECEIVE, Timestamp: time.Now().Unix()}
}

func (tx *Transaction) SetHash() error {
	hash, err := tx.GetHash()
	if err != nil {
		return err
	}
	tx.Hash = hash
	return nil
}

func (tx *Transaction) GetHash() ([]byte, error) {
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
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, tx.Hash)
	if err != nil {
		return err
	}
	tx.Signature = append(r.Bytes(), s.Bytes()...)
	return nil
}

func (tx *Transaction) VerifySignature() bool {
	r := big.Int{}
	s := big.Int{}
	sigLen := len(tx.Signature)
	r.SetBytes(tx.Signature[:(sigLen / 2)])
	s.SetBytes(tx.Signature[(sigLen / 2):])

	x := big.Int{}
	y := big.Int{}
	keyLen := len(tx.PubKey)
	x.SetBytes(tx.PubKey[:(keyLen / 2)])
	y.SetBytes(tx.PubKey[(keyLen / 2):])

	curve := elliptic.P256()
	rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}

	return ecdsa.Verify(&rawPubKey, tx.Hash, &r, &s)
}

func (tx *Transaction) ToBytes() []byte {
	var result bytes.Buffer
	encoder := xdr.NewEncoder(&result)
	encoder.Encode(tx)
	return result.Bytes()
}

func NewTransactionFromBytes(d []byte) *Transaction {
	var tx Transaction
	decoder := xdr.NewDecoder(bytes.NewReader(d))
	decoder.Decode(&tx)
	return &tx
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
