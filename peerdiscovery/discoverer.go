package peerdiscovery

import "github.com/msaldanha/realChain/consensus"

//go:generate mockgen -destination=../tests/mock_discoverer.go -package=tests github.com/msaldanha/realChain/peerdiscovery Discoverer

type Discoverer interface {
	Init() error
	Peers() ([]consensus.ConsensusClient, error)
}
