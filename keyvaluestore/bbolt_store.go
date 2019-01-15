package keyvaluestore

import (
	"github.com/coreos/bbolt"
	"github.com/msaldanha/realChain/errors"
	"time"
)

const (
	ErrExpectedBoltKeyValueStoreOptions = errors.Error("expected BoltKeyValueStoreOptions type")
	ErrInvalidBucketName                = errors.Error("invalid bucket name")
)

type BoltKeyValueStoreOptions struct {
	BucketName string
	DbFile     string
}

type BoltKeyValueStore struct {
	db         *bolt.DB
	BucketName string
}

func NewBoltKeyValueStore() *BoltKeyValueStore {
	return &BoltKeyValueStore{}
}

func (st *BoltKeyValueStore) Init(options interface{}) error {
	if _, ok := options.(*BoltKeyValueStoreOptions); !ok {
		return ErrExpectedBoltKeyValueStoreOptions
	}

	opt := options.(*BoltKeyValueStoreOptions)
	if opt.BucketName == "" {
		return ErrInvalidBucketName
	}

	db, err := bolt.Open(opt.DbFile, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(opt.BucketName))
		if err != nil {
			return err
		}
		return nil
	})
	st.db = db
	st.BucketName = opt.BucketName
	return nil
}

func (st *BoltKeyValueStore) Put(key string, value []byte) (error) {
	return st.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(st.BucketName))
		err := b.Put([]byte(key), value)
		return err
	})
}

func (st *BoltKeyValueStore) Get(key string) (ret []byte, ok bool, err error) {
	ok = false
	ret = nil
	err = st.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(st.BucketName))
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
		b := tx.Bucket([]byte(st.BucketName))
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
		b := tx.Bucket([]byte(st.BucketName))
		c := b.Cursor()
		size = 0
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			size = size + 1
		}

		return nil
	})
	return
}

func (st *BoltKeyValueStore) GetAll() ([][]byte, error) {
	all := make([][]byte, 0)
	err := st.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(st.BucketName))
		b.ForEach(func(k, v []byte) error {
			ret := make([]byte, len(v))
			copy(ret, v)
			all = append(all, ret)
			return nil
		})
		return nil
	})
	return all, err
}
