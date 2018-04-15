package transactionstore

import (
	"github.com/msaldanha/realChain/transaction"
	"github.com/msaldanha/realChain/keyvaluestore"
	"math/big"
	"bytes"
	"math"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"encoding/hex"
	"time"
)

const targetBits int16 = 16

type TransactionStore struct {
	store            keyvaluestore.Storer
	validatorCreator transaction.ValidatorCreator
}

func New(store keyvaluestore.Storer, validatorCreator transaction.ValidatorCreator) (*TransactionStore) {
	a := &TransactionStore{store: store, validatorCreator: validatorCreator}
	return a
}

func (ts *TransactionStore) isValid(tx *transaction.Transaction) (bool, error) {
	if !tx.Type.IsValid(){
		return false, transaction.ErrInvalidTransactionType
	}
	val := ts.validatorCreator.CreateValidatorForTransaction(tx.Type, ts.store)
	return val.IsValid(tx)
}

func (ts *TransactionStore) Store(tx *transaction.Transaction) (*transaction.Transaction, error) {
	if ok, err := ts.isValid(tx); !ok {
		return nil, err
	}
	ts.store.Put(string(tx.Hash), tx.ToBytes())
	ts.store.Put(string(tx.Address), tx.ToBytes())
	return tx, nil
}

func (ts *TransactionStore) Retrieve(hash string) (*transaction.Transaction, error) {
	value, _, err := ts.GetTransaction(hash)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (ts *TransactionStore) CalculatePow(tx *transaction.Transaction) (int64, []byte, error) {
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

func (ts *TransactionStore) VerifyPow(tx *transaction.Transaction) (bool, error) {
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

func (ts *TransactionStore) CreateOpenTransaction() (*transaction.Transaction) {
	return &transaction.Transaction{Type: transaction.OPEN, Timestamp: time.Now().Unix()}
}

func (ts *TransactionStore) CreateSendTransaction() (*transaction.Transaction) {
	return &transaction.Transaction{Type: transaction.SEND, Timestamp: time.Now().Unix()}
}

func (ts *TransactionStore) CreateReceiveTransaction() (*transaction.Transaction) {
	return &transaction.Transaction{Type: transaction.RECEIVE, Timestamp: time.Now().Unix()}
}

func (ts *TransactionStore) GetTransactionChain(txHash string, includeAll bool) ([]*transaction.Transaction, error) {
	tx, ok, _ := ts.GetTransaction(txHash)
	chain := []*transaction.Transaction{}
	for ok {
		chain = append(chain[:0], append([]*transaction.Transaction{tx}, chain[0:]...)...)
		if len(tx.Previous) > 0 {
			tx, ok, _ = ts.GetTransaction(string(tx.Previous))
		} else if tx.Type == transaction.OPEN && len(tx.Link) > 0 && includeAll {
			tx, ok, _ = ts.GetTransaction(string(tx.Link))
		} else {
			break
		}
	}
	return chain, nil
}


func (ts *TransactionStore) GetTransaction(txHash string) (*transaction.Transaction, bool, error) {
	tx, ok, err := ts.store.Get(txHash)
	if tx == nil {
		return nil, ok, err
	}
	return transaction.NewTransactionFromBytes(tx), ok, err
}

func (ts *TransactionStore) IsEmpty() (bool) {
	return ts.store.IsEmpty()
}

func getTarget() (*big.Int) {
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

