package keyvaluestore

type MemoryKeyValueStore struct {
	pairs map[string][]byte
	tip []byte
}

func NewMemoryKeyValueStore() (*MemoryKeyValueStore) {
	pairs := make(map[string][]byte)
	return &MemoryKeyValueStore{pairs:pairs}
}

func (st *MemoryKeyValueStore) Init(options interface{}) (error) {
	return nil
}

func (st *MemoryKeyValueStore) Put(key string, value []byte) (error) {
	st.tip = value
	st.pairs[key] = value
	return nil
}

func (st *MemoryKeyValueStore) Get(key string) ([]byte, bool, error) {
	value, found := st.pairs[key]
	return value, found, nil
}

func (st *MemoryKeyValueStore) GetTip(key string) ([]byte, bool, error) {
	if st.tip == nil {
		return nil, false, nil
	}
	return st.tip, true, nil
}

func (st *MemoryKeyValueStore) IsEmpty() (bool) {
	return len(st.pairs) == 0
}

func (st *MemoryKeyValueStore) Size() (int) {
	return len(st.pairs)
}