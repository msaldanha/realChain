package peerdiscovery

import "github.com/msaldanha/realChain/consensus"

type StaticDiscoverer struct {

}

func NewStaticDiscoverer() *StaticDiscoverer{
	return &StaticDiscoverer{}
}

func (d *StaticDiscoverer) Init() error {
	return nil
}

func (d *StaticDiscoverer) Peers() ([]consensus.ConsensusClient, error) {
	peers := [0]consensus.ConsensusClient{}
	return peers[:], nil
}
