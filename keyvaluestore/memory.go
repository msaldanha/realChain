package keyvaluestore

import "github.com/msaldanha/realChain/block"

type MemoryKeyValueStore struct {
	pairs map[string]*block.Block
	tip *block.Block
}

func NewMemoryKeyValueStore() (*MemoryKeyValueStore) {
	pairs := make(map[string]*block.Block)
	return &MemoryKeyValueStore{pairs:pairs}
}

func (st *MemoryKeyValueStore) Put(key string, value *block.Block) (error) {
	st.tip = value
	st.pairs[key] = value
	return nil
}

func (st *MemoryKeyValueStore) Get(key string) (*block.Block, bool, error) {
	value, found := st.pairs[key]
	return value, found, nil
}

func (st *MemoryKeyValueStore) GetTip(key string) (*block.Block, bool, error) {
	if st.tip == nil {
		return nil, false, nil
	}
	return st.tip, true, nil
}

func (st *MemoryKeyValueStore) IsEmpty() (bool) {
	return len(st.pairs) == 0
}