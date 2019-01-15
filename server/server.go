package server

import (
	"github.com/msaldanha/realChain/errors"
	"github.com/msaldanha/realChain/consensus"
	"github.com/msaldanha/realChain/ledger"
	"github.com/msaldanha/realChain/peerdiscovery"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
)

//go:generate mockgen -destination=../tests/mock_listener.go -package=tests net Listener

const (
	ErrDeclinedByVoting                 = errors.Error("transaction declined by voting")
	ErrNoPeersForVoting                 = errors.Error("no peers for voting")
)

type Server struct {
	ledger    ledger.Ledger
	consensus consensus.Consensus
	discovery peerdiscovery.Discoverer
	listener  net.Listener
}

func New(ledger ledger.Ledger,
		consensus consensus.Consensus,
		discovery peerdiscovery.Discoverer,
		listener net.Listener) *Server {
	return &Server{ledger: ledger, consensus: consensus, discovery: discovery, listener: listener}
}

func (s *Server) Run() error {
	grpcServer := grpc.NewServer()
	consensus.RegisterConsensusServer(grpcServer, s)
	ledger.RegisterLedgerServer(grpcServer, s)
	return grpcServer.Serve(s.listener)
}

func (s *Server) Register(ctx context.Context, request *ledger.RegisterRequest) (*ledger.RegisterResult, error) {
	err := s.ledger.Verify(request.SendTx, request.ReceiveTx)
	if err != nil {
		return nil, err
	}

	err = s.resolve(ctx, request)
	if err != nil {
		return nil, err
	}
	return &ledger.RegisterResult{}, nil
}

func (s *Server) GetLastTransaction(ctx context.Context, request *ledger.GetLastTransactionRequest) (*ledger.GetLastTransactionResult, error) {
	tx, err := s.ledger.GetLastTransaction(request.Address)
	if err != nil {
		return nil, err
	}
	return &ledger.GetLastTransactionResult{Tx: tx}, nil
}

func (s *Server) GetTransaction(ctx context.Context, request *ledger.GetTransactionRequest) (*ledger.GetTransactionResult, error) {
	tx, err := s.ledger.GetTransaction(request.Hash)
	if err != nil {
		return nil, err
	}
	return &ledger.GetTransactionResult{Tx: tx}, nil
}

func (s *Server) VerifyTransaction(ctx context.Context, request *ledger.VerifyTransactionRequest) (*ledger.VerifyTransactionResult, error) {
	err := s.ledger.VerifyTransaction(request.Tx, true)
	if err != nil {
		return nil, err
	}
	return &ledger.VerifyTransactionResult{}, nil
}

func (s *Server) Verify(ctx context.Context, request *ledger.VerifyRequest) (*ledger.VerifyResult, error) {
	err := s.ledger.Verify(request.SendTx, request.ReceiveTx)
	if err != nil {
		return nil, err
	}
	return &ledger.VerifyResult{}, nil
}

func (s *Server) Vote(ctx context.Context, request *consensus.VoteRequest) (*consensus.VoteResult, error) {
	return s.consensus.Vote(request)
}

func (s *Server) Accept(ctx context.Context, request *consensus.AcceptRequest) (*consensus.AcceptResult, error) {
	return s.consensus.Accept(request)
}

func (s *Server) resolve(ctx context.Context, request *ledger.RegisterRequest) error {
	peers, err := s.discovery.Peers()
	if err != nil {
		return err
	}

	if len(peers) == 0 {
		return ErrNoPeersForVoting
	}

	nok := 0
	votes := []*consensus.Vote{}
	for _, peer := range peers {
		result, err := peer.Vote(ctx, &consensus.VoteRequest{SendTx: request.SendTx, ReceiveTx: request.ReceiveTx})
		if err != nil {
			return err
		}
		if !result.Vote.Ok {
			nok++
		}
		votes = append(votes, result.Vote)
	}

	if nok > 0 {
		return ErrDeclinedByVoting
	}

	err = s.ledger.Register(request.SendTx, request.ReceiveTx)
	if err != nil {
		return err
	}

	accept := &consensus.AcceptRequest{SendTx: request.SendTx, ReceiveTx: request.ReceiveTx, Votes: votes}
	for _, peer := range peers {
		_, err = peer.Accept(ctx, accept)
	}
	return nil
}
