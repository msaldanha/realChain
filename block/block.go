package block

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
