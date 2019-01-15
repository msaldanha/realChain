package consensus

import (
	"github.com/msaldanha/realChain/address"
	"github.com/msaldanha/realChain/errors"
	"github.com/msaldanha/realChain/ledger"
)

//go:generate mockgen -destination=../tests/mock_consensusclient.go -package=tests github.com/msaldanha/realChain/consensus ConsensusClient
//go:generate mockgen -destination=../tests/mock_consensus.go -package=tests github.com/msaldanha/realChain/consensus Consensus
//go:generate protoc -I.. consensus/consensus.proto --go_out=plugins=grpc:../

const (
	ErrInvalidVotingResult = errors.Error("invalid voting result")
)

type Consensus interface {
	Vote(*VoteRequest) (*VoteResult, error)
	Accept(*AcceptRequest) (*AcceptResult, error)
}

type consensus struct {
	ledger ledger.Ledger
	address *address.Address
}

func NewConsensus(ledger ledger.Ledger, address *address.Address) *consensus {
	return &consensus{
		ledger: ledger,
		address: address,
	}
}

func (c *consensus) Vote(request *VoteRequest) (*VoteResult, error) {
	err := c.ledger.VerifyTransaction(request.ReceiveTx, true)
	if err != nil {
		return c.createResult(false, err.Error())
	}

	err = c.ledger.VerifyTransaction(request.SendTx, true)
	if err != nil {
		return c.createResult(false, err.Error())
	}

	err = c.ledger.Verify(request.SendTx, request.ReceiveTx)
	if err != nil {
		return c.createResult(false, err.Error())
	}

	return c.createResult(true, "")
}

func (c *consensus) Accept(request *AcceptRequest) (*AcceptResult, error) {
	oks := 0
	for _, vote := range request.Votes {
		err := c.validateVote(vote)
		if err != nil {
			return nil, err
		}
		if vote.Ok {
			oks++
		}
	}
	if oks == 0 || oks != len(request.Votes) {
		return nil, ErrInvalidVotingResult
	}
	err := c.ledger.Register(request.SendTx, request.ReceiveTx)
	if err != nil {
		return nil, err
	}
	return &AcceptResult{}, nil
}

func (c *consensus) createResult(ok bool, reason string) (*VoteResult, error) {
	vote := &Vote{Ok: ok, Reason: reason}
	err := vote.Sign(c.address.Keys.ToEcdsaPrivateKey())
	if err != nil {
		return nil, err
	}
	vote.PubKey = c.address.Keys.PublicKey
	return &VoteResult{Vote: vote}, nil
}

func (c *consensus) validateVote(vote *Vote) error {
	return nil
}

