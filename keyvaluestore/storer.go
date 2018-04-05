package keyvaluestore

type Storer interface {
	Put(key string, value interface{}) (error)
	Get(key string) (interface{}, bool, error)
	GetTip(key string) (interface{}, bool, error)
	IsEmpty() (bool)
}
