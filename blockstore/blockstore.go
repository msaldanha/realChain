package blockstore

import (
	"github.com/msaldanha/realChain/block"
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

type BlockStore struct {
	store                 keyvaluestore.Storer
	blockValidatorCreator block.BlockValidatorCreator
}

func New(store keyvaluestore.Storer, validatorCreator block.BlockValidatorCreator) (*BlockStore) {
	a := &BlockStore{store: store, blockValidatorCreator: validatorCreator}
	return a
}

func (bs *BlockStore) isValid(blk *block.Block) (bool, error) {
	if !blk.Type.IsValid(){
		return false, block.ErrInvalidBlockType
	}
	val := bs.blockValidatorCreator.CreateValidatorForBlock(blk.Type, bs.store)
	return val.IsValid(blk)
}

func (bs *BlockStore) Store(blk *block.Block) (*block.Block, error) {
	if ok, err := bs.isValid(blk); !ok {
		return nil, err
	}
	bs.store.Put(string(blk.Hash), blk)
	bs.store.Put(string(blk.Account), blk)
	return blk, nil
}

func (bs *BlockStore) Retrieve(hash string) (*block.Block, error) {
	value, _, err := bs.GetBlock(hash)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (bs *BlockStore) CalculatePow(blk *block.Block) (int64, []byte, error) {
	var hashInt big.Int
	var hash [32]byte
	var nonce int64 = 0

	target := getTarget()

	data, err := blk.GetHashableBytes()
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

func (bs *BlockStore) VerifyPow(blk *block.Block) (bool, error) {
	var hashInt big.Int

	target := getTarget()

	data, err := blk.GetHashableBytes()
	if err != nil {
		return false, err
	}
	dataWithNonce := append(data, int64ToBytes(blk.PowNonce))
	hash := sha256.Sum256(bytes.Join(dataWithNonce, []byte{}))
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(target) == -1, nil
}

func (bs *BlockStore) CreateOpenBlock() (*block.Block) {
	return &block.Block{Type:block.OPEN, Timestamp: time.Now().Unix()}
}

func (bs *BlockStore) CreateSendBlock() (*block.Block) {
	return &block.Block{Type:block.SEND, Timestamp: time.Now().Unix()}
}

func (bs *BlockStore) CreateReceiveBlock() (*block.Block) {
	return &block.Block{Type:block.RECEIVE, Timestamp: time.Now().Unix()}
}

func (bs *BlockStore) GetBlockChain(blockHash string) ([]*block.Block, error) {
	blk, ok, _ := bs.GetBlock(blockHash)
	chain := []*block.Block{}
	for ok {
		chain = append(chain[:0], append([]*block.Block{blk}, chain[0:]...)...)
		if len(blk.Previous) > 0 {
			blk, ok, _ = bs.GetBlock(string(blk.Previous))
		} else if blk.Type == block.OPEN && len(blk.Link) > 0 {
			blk, ok, _ = bs.GetBlock(string(blk.Link))
		} else {
			break
		}
	}
	return chain, nil
}


func (bs *BlockStore) GetBlock(blockHash string) (*block.Block, bool, error) {
	blk, ok, err := bs.store.Get(blockHash)
	if blk == nil {
		return nil, ok, err
	}
	return blk.(*block.Block), ok, err
}

func (bs *BlockStore) IsEmpty() (bool) {
	return bs.store.IsEmpty()
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

