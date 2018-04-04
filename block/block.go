package block

type BlockType int

const (
	OPEN BlockType = 1 << iota
	SEND
	RECEIVE
	CHANGE
)

type Block struct {
	Type BlockType
	Account string
	Representative string
	Previous string
	Source string
	Destination string
	Work string
	Balance float64
	Signature string
}

func (bt BlockType) IsValid() (bool) {
	return bt >= OPEN && bt <= CHANGE
}
