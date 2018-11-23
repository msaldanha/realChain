package consensus

import (
	"github.com/msaldanha/realChain/network"
	"github.com/msaldanha/realChain/ledger"
	"github.com/msaldanha/realChain/transaction"
	"log"
	"bytes"
	"github.com/davecgh/go-xdr/xdr2"
	"github.com/msaldanha/realChain/Error"
	"time"
)

const (
	VoteOk                   int8  = 1
	VoteNok                  int8  = -1
	Magic                    uint32 = 0xfffecaba
	TransactionVotingTimeout       = time.Second * 10

	EndPointVoting      = "consensus.v1.voting"
	EndPointVote        = "consensus.v1.vote"
	EndPointTransaction = "consensus.v1.transaction"

	ErrNoNetwork = Error.Error("no network")
)

type Consensus struct {
	network             *network.Network
	ledger              ledger.Ledger
	txPool              map[string]*transaction.Transaction
	currentVoting       *transaction.Transaction
	toVoteQueue         chan *voting
	finishedVotingQueue chan *voting
	votingContext       *voting
}

type vote struct {
	txHash string
	value  int8
	reason string
}

type voting struct {
	tx       *transaction.Transaction
	votes    map[string]*vote
	ok       int
	nok      int
	failed   int
	expected int
	err      error
}

func New(network *network.Network, ledger ledger.Ledger) *Consensus {
	return &Consensus{
		network:             network,
		ledger:              ledger,
		txPool:              make(map[string]*transaction.Transaction),
		toVoteQueue:         make(chan *voting),
		finishedVotingQueue: make(chan *voting),
	}
}

func (v *vote) ToBytes() []byte {
	var result bytes.Buffer
	encoder := xdr.NewEncoder(&result)
	encoder.Encode(v)
	return result.Bytes()
}

func (c *Consensus) Run() error {
	if c.network == nil {
		return ErrNoNetwork
	}
	c.network.InstallHandler(EndPointVoting, c.handleVoting())
	c.network.InstallHandler(EndPointVote, c.handleVote())
	c.network.InstallHandler(EndPointTransaction, c.handleTransaction())
	c.doVoting()
	return nil
}

func (c *Consensus) handleTransaction() network.Handler {
	return func(ctx *network.Context) {
		data := ctx.Data
		tx := transaction.NewTransactionFromBytes(data)

		if txInPool := c.txPool[string(tx.Hash)]; txInPool != nil {
			return
		}

		if err := c.ledger.VerifyTransaction(tx, true); err != nil {
			logError("handle", tx, err)
			return
		}

		c.txPool[string(tx.Hash)] = tx
		c.network.Broadcast(EndPointTransaction, data)
		c.toVoteQueue <- &voting{tx: tx, votes: make(map[string]*vote)}
	}
}

func (c *Consensus) handleVote() network.Handler {
	return func(ctx *network.Context) {
		var err error
		data := ctx.Data
		vote := &vote{txHash: string(data), value: VoteNok}

		tx := c.txPool[vote.txHash]
		if tx == nil {
			tx, err = c.ledger.GetLastTransaction(vote.txHash)
		}

		if tx != nil && err == nil {
			vote.value = VoteOk
			ctx.Peer.Send(EndPointVoting, vote.ToBytes())
			return
		}

		if tx == nil && err == nil {
			vote.value = VoteNok
			vote.reason = "Unknown transaction"
			ctx.Peer.Send(EndPointVoting, vote.ToBytes())
			return
		}

		if err != nil {
			vote.reason = err.Error()
		}

		ctx.Peer.Send(EndPointVoting, vote.ToBytes())
	}
}

func (c *Consensus) handleVoting() network.Handler {
	return func(ctx *network.Context) {
		data := ctx.Data
		vote, _ := newVoteFromBytes(data)
		if vote == nil {
			return
		}

		peerVote := c.votingContext.votes[ctx.Peer.String()]
		if peerVote != nil {
			return
		}

		c.votingContext.votes[ctx.Peer.String()] = vote

		if vote.value == VoteOk {
			c.votingContext.ok++
		} else if vote.value == VoteNok {
			c.votingContext.nok++
		} else {
			c.votingContext.failed++
		}

		if len(c.votingContext.votes) == c.votingContext.expected {
			c.finishedVotingQueue <- c.votingContext
		}
	}
}

func (c *Consensus) confirmTransaction(tx *transaction.Transaction) () {
	//TODO: fix me
	//c.ledger.Register(tx)
	delete(c.txPool, string(tx.Hash))
}

func (c *Consensus) cleanupTransaction(tx *transaction.Transaction) () {
	delete(c.txPool, string(tx.Hash))
}

func (c *Consensus) doVoting() {
	for {
		voting := <-c.toVoteQueue
		//is there a voting in progress ?
		if c.votingContext != nil {
			panic("Should handle only one voting at a time!")
		}
		c.votingContext.expected = c.network.GetPeersCount()
		c.network.Broadcast(EndPointVote, voting.tx.Hash)
		//wait for at most TransactionVotingTimeout to have the voting finished
		c.votingContext = voting
		timer := time.NewTimer(TransactionVotingTimeout)
		select {
		case voting = <-c.finishedVotingQueue:
			if voting.err == nil && votingPassed(voting) {
				c.confirmTransaction(voting.tx)
			}
		case <-timer.C:
			c.confirmTransaction(voting.tx)
			break
		}
		timer.Stop()
		c.votingContext = nil
	}
}

func votingPassed(vo *voting) bool {
	total := vo.ok + vo.nok + vo.failed
	return (vo.ok / total * 100) >= 66
}

func newVoteFromBytes(d []byte) (*vote, error) {
	var tx vote
	decoder := xdr.NewDecoder(bytes.NewReader(d))
	_, err := decoder.Decode(&tx)
	return &tx, err
}

func logError(context string, tx *transaction.Transaction, err error) {
	log.Printf("Consensus.%s failed: %s (tx: %s)", context, err, string(tx.Hash))
}
