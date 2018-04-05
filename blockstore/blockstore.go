package blockstore

import (
	. "github.com/msaldanha/realChain/block"
	. "github.com/msaldanha/realChain/keyvaluestore"
	"errors"
	"github.com/msaldanha/realChain/validator"
	"math/big"
	"bytes"
	"math"
	"crypto/sha256"
	"fmt"
	"encoding/binary"
	"log"
	"encoding/hex"
	"time"
)

const targetBits int16 = 16

type BlockStore struct {
	store                 Storer
	blockValidatorCreator validator.BlockValidatorCreator
}

func New(store Storer, validatorCreator validator.BlockValidatorCreator) (*BlockStore) {
	a := &BlockStore{store: store, blockValidatorCreator: validatorCreator}
	return a
}

func (bs *BlockStore) isValid(block *Block) (bool, error) {
	if !block.Type.IsValid(){
		return false, errors.New("Invalid block type")
	}
	val := bs.blockValidatorCreator.CreateValidatorForBlock(block.Type, bs.store)
	return val.IsValid(block)
}

func (bs *BlockStore) Store(block *Block) (*Block, error) {
	if ok, err := bs.isValid(block); !ok {
		return nil, err
	}
	bs.store.Put(string(block.Hash), block)
	return block, nil
}

func (bs *BlockStore) Retrieve(hash string) (*Block, error) {
	value, _, err := bs.store.Get(hash)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (bs *BlockStore) CalculatePow(block *Block) (int64, []byte, error) {
	var hashInt big.Int
	var hash [32]byte
	var nonce int64 = 0

	target := getTarget()

	data, err := block.GetHashableBytes()
	if err != nil {
		return 0, nil, err
	}

	for nonce < math.MaxInt64 {
		dataWithNonce := append(data, int64ToBytes(nonce))
		hash = sha256.Sum256(bytes.Join(dataWithNonce, []byte{}))
		//fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")

	hexHash := []byte(hex.EncodeToString(hash[:]))

	return nonce, hexHash[:], nil
}

func (bs *BlockStore) VerifyPow(block *Block) (bool, error) {
	var hashInt big.Int

	target := getTarget()

	data, err := block.GetHashableBytes()
	if err != nil {
		return false, err
	}
	dataWithNonce := append(data, int64ToBytes(block.PowNonce))
	hash := sha256.Sum256(bytes.Join(dataWithNonce, []byte{}))
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(target) == -1, nil
}

func (bs *BlockStore) CreateOpenBlock() (*Block) {
	return &Block{Type:OPEN, Timestamp: time.Now().Unix()}
}

func (bs *BlockStore) CreateSendBlock() (*Block) {
	return &Block{Type:SEND, Timestamp: time.Now().Unix()}
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

