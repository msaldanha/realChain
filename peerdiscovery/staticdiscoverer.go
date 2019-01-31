package peerdiscovery

import (
	"github.com/msaldanha/realChain/config"
	"github.com/msaldanha/realChain/consensus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

type StaticDiscoverer struct {
	cfg *viper.Viper
	peers []consensus.ConsensusClient
}

func NewStaticDiscoverer(cfg *viper.Viper) *StaticDiscoverer{
	return &StaticDiscoverer{cfg: cfg}
}

func (d *StaticDiscoverer) Init() error {
	return nil
}

func (d *StaticDiscoverer) Peers() ([]consensus.ConsensusClient, error) {
	if d.peers != nil {
		return d.peers, nil
	}

	peers := make([]consensus.ConsensusClient, 0)
	ips := d.cfg.GetStringSlice(config.CfgPeers)
	for _, v := range ips {
		conn, err := grpc.Dial(v, grpc.WithInsecure())
		if err == nil {
			peer := consensus.NewConsensusClient(conn)
			peers = append(peers, peer)
		}
	}

	d.peers = peers
	return d.peers, nil
}
