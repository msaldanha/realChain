package keyvaluestore

import "github.com/msaldanha/realChain/block"

type Storer interface {
	Put(key string, value *block.Block) (error)
	Get(key string) (*block.Block, bool, error)
	GetTip(key string) (*block.Block, bool, error)
	IsEmpty() (bool)
	Size() (int)
}
