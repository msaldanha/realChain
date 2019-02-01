package ledger

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"github.com/golang/protobuf/proto"
	"github.com/msaldanha/realChain/address"
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

func CreateGenesisTransaction(balance float64) (*Transaction, *address.Address, error) {
	genesisTx := NewOpenTransaction()
	addr, err := address.NewAddressWithKeys()
	if err != nil {
		return nil, nil, err
	}

	genesisTx.Address = addr.Address
	genesisTx.Balance = balance
	genesisTx.PubKey = hex.EncodeToString(addr.Keys.PublicKey)

	err = genesisTx.SetPow()
	if err != nil {
		return nil, nil, err
	}

	err = genesisTx.Sign(addr.Keys.ToEcdsaPrivateKey())
	if err != nil {
		return nil, nil, err
	}

	return genesisTx, addr, nil
}

func CreateSendTransaction(fromTipTx *Transaction, fromAddr *address.Address, to string,
		amount float64) (*Transaction, error) {
	sendTx := NewSendTransaction()
	sendTx.Address = fromTipTx.Address
	sendTx.Link = to
	sendTx.Previous = fromTipTx.Hash
	sendTx.Balance = fromTipTx.Balance - amount
	sendTx.PubKey = fromTipTx.PubKey

	if err := sendTx.SetPow(); err != nil {
		return nil, err
	}

	if err := sendTx.Sign(fromAddr.Keys.ToEcdsaPrivateKey()); err != nil {
		return nil, err
	}

	return sendTx, nil
}

func CreateReceiveTransaction(send *Transaction, amount float64, receiveAddr *address.Address,
		receiveTipTx *Transaction) (*Transaction, error) {
	var receiveTx *Transaction
	if receiveTipTx != nil {
		receiveTx = NewReceiveTransaction()
		receiveTx.Previous = receiveTipTx.Hash
		receiveTx.Balance = receiveTipTx.Balance + amount
		receiveTx.PubKey = receiveTipTx.PubKey
	} else {
		receiveTx = NewOpenTransaction()
		receiveTx.Balance = amount
		receiveTx.PubKey = hex.EncodeToString(receiveAddr.Keys.PublicKey)
	}

	receiveTx.Address = send.Link
	receiveTx.Link = send.Hash

	if err := receiveTx.SetPow(); err != nil {
		return nil, err
	}

	if err := receiveTx.Sign(receiveAddr.Keys.ToEcdsaPrivateKey()); err != nil {
		return nil, err
	}

	return receiveTx, nil
}

func (tx *Transaction) SetHash() error {
	hash, err := tx.CalculateHash()
	if err != nil {
		return err
	}
	tx.Hash = hash
	return nil
}

func (tx *Transaction) CalculateHash() (string, error) {
	hashableBytes, err := tx.GetHashableBytes()
	if err != nil {
		return "", err
	}
	headers := bytes.Join(hashableBytes, []byte{})
	hash := sha256.Sum256(headers)
	return hex.EncodeToString(hash[:]), nil
}

func (tx *Transaction) GetHashableBytes() ([][]byte, error) {
	var balance bytes.Buffer
	if err := binary.Write(&balance, binary.LittleEndian, tx.Balance); err != nil {
		return nil, err
	}

	timestamp := []byte(strconv.FormatInt(tx.Timestamp, 10))

	return [][]byte{timestamp, []byte(tx.Address), []byte(tx.Previous), []byte(tx.Link), balance.Bytes()}, nil
}

func (tx *Transaction) CalculatePow() (int64, string, error) {
	var hashInt big.Int
	var hash [32]byte
	var nonce int64 = 0

	target := getTarget()

	data, err := tx.GetHashableBytes()
	if err != nil {
		return 0, "", err
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

	hexHash := hex.EncodeToString(hash[:])

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
	hash, err := hex.DecodeString(tx.Hash)
	if err != nil {
		return err
	}

	s, err := crypto.Sign(hash, privateKey)
	if err != nil {
		return err
	}

	tx.Signature = hex.EncodeToString(s)
	return nil
}

func (tx *Transaction) VerifySignature() bool {
	sign, err := hex.DecodeString(tx.Signature)
	if err != nil {
		return false
	}

	pubKey, err := hex.DecodeString(tx.PubKey)
	if err != nil {
		return false
	}

	hash, err := hex.DecodeString(tx.Hash)
	if err != nil {
		return false
	}

	return crypto.VerifySignature(sign, pubKey, hash)
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
