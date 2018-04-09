package ledger

import (
	"github.com/msaldanha/realChain/keypair"
)

type Account struct {
	Keys *keypair.KeyPair
	Address string
}
