package keyvaluestore

type Storer interface {
	Init(options interface{}) (error)
	Put(key string, value []byte) (error)
	Get(key string) ([]byte, bool, error)
	GetAll() ([][]byte, error)
	GetTip(key string) ([]byte, bool, error)
	IsEmpty() (bool)
	Size() (int)
}
