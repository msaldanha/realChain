package keyvaluestore

type MemoryKeyValueStore struct {
	pairs map[string]interface{}
	tip interface{}
}

func NewMemoryKeyValueStore() (*MemoryKeyValueStore) {
	pairs := make(map[string]interface{})
	return &MemoryKeyValueStore{pairs:pairs}
}

func (st *MemoryKeyValueStore) Put(key string, value interface{}) (error) {
	st.tip = value
	st.pairs[key] = value
	return nil
}

func (st *MemoryKeyValueStore) Get(key string) (interface{}, bool, error) {
	value, found := st.pairs[key]
	return value, found, nil
}

func (st *MemoryKeyValueStore) GetTip(key string) (interface{}, bool, error) {
	if st.tip == nil {
		return nil, false, nil
	}
	return st.tip, true, nil
}