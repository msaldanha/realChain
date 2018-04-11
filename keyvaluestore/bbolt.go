package keyvaluestore

import (
	"github.com/coreos/bbolt"
	"github.com/msaldanha/realChain/Error"
	"time"
)

const (
	ErrExpectedBoltKeyValueStoreOptions               = Error.Error("expected BoltKeyValueStoreOptions type")

	bucketName = "BlockChain"
)


type BoltKeyValueStoreOptions struct {
	DbFile string
}

type BoltKeyValueStore struct {
	db *bolt.DB
	blockchain *bolt.Bucket
}

func NewBoltKeyValueStore() (*BoltKeyValueStore) {
	return &BoltKeyValueStore{}
}

func (st *BoltKeyValueStore) Init(options interface{}) (error) {
	if _, ok := options.(*BoltKeyValueStoreOptions); !ok {
		return ErrExpectedBoltKeyValueStoreOptions
	}

	opt := options.(*BoltKeyValueStoreOptions)
	db, err := bolt.Open(opt.DbFile, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}

	var blockchain *bolt.Bucket
	db.Update(func(tx *bolt.Tx) error {
		blockchain, err = tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}
		return nil
	})
	st.db = db
	st.blockchain = blockchain
	return nil
}

func (st *BoltKeyValueStore) Put(key string, value []byte) (error) {
	return st.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		err := b.Put([]byte(key), value)
		return err
	})
}

func (st *BoltKeyValueStore) Get(key string) (ret []byte, ok bool, err error) {
	ok = false
	ret = nil
	err = st.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		value := b.Get([]byte(key))
		if len(value) == 0 {
			return nil
		}
		ok = true
		ret = make([]byte, len(value))
		copy(ret, value)
		return nil
	})
	return
}

func (st *BoltKeyValueStore) GetTip(key string) ([]byte, bool, error) {
	return nil, false, nil
}

func (st *BoltKeyValueStore) IsEmpty() (isEmpty bool) {
	isEmpty = true
	st.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		c := b.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			isEmpty = false
			return nil
		}
		return nil
	})
	return
}

func (st *BoltKeyValueStore) Size() (size int) {
	st.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		c := b.Cursor()
		size = 0
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			size = size + 1
		}

		return nil
	})
	return
}
