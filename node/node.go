package node

import (
	"github.com/msaldanha/realChain/ledger"
	"github.com/msaldanha/realChain/keyvaluestore"
	"github.com/msaldanha/realChain/block"
	"github.com/msaldanha/realChain/blockstore"
	"fmt"
	"os"
)

type Node struct {
	ld *ledger.Ledger
}

func New() (*Node) {
	return &Node{}
}

func (n *Node) Run() {

	ld, err := createLedger()
	checkError(err)

	udp := NewUdpServer(ld)
	rest := NewRestServer(ld)

	udpch := make(chan error)
	restch := make(chan error)

	go func() {udpch <- udp.Run()}()
	go func() {restch <- rest.Run()}()

	for {
		select {
		case eu := <- udpch:
			checkError(eu)
		case er := <- restch:
			checkError(er)
		}
	}
	fmt.Println("Done.")
}

func checkError(err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}

func createLedger() (*ledger.Ledger, error) {
	ms := keyvaluestore.NewMemoryKeyValueStore()
	as := keyvaluestore.NewMemoryKeyValueStore()
	val := block.NewBlockValidatorCreator()
	bs := blockstore.New(ms, val)
	ld := ledger.New()
	ld.Use(bs, as)
	if bs.IsEmpty() {
		if _, err := ld.Initialize(10000); err != nil {
			return nil, err
		}
	}
	return ld, nil
}