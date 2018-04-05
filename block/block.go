package block

import "fmt"

type BlockType int

const (
	OPEN BlockType = 1 + iota
	SEND
	RECEIVE
	CHANGE
)

type Block struct {
	Type BlockType
	Account string
	Representative string
	Previous string
	Link string
	Work string
	Balance float64
	Hash string
	Signature string
}

func (bt BlockType) IsValid() (bool) {
	return bt >= OPEN && bt <= CHANGE
}

func (bt BlockType) String() (string) {
	name := ""
	switch bt {
	case OPEN:
		name = "OPEN"
	case SEND:
		name = "SEND"
	case RECEIVE:
		name = "RECEIVE"
	case CHANGE:
		name = "CHANGE"
	}
	return fmt.Sprintf("%d(%s)", int(bt), name)
}
