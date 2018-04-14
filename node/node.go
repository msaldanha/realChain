package node

import (
	"github.com/msaldanha/realChain/ledger"
	"github.com/msaldanha/realChain/keyvaluestore"
	"github.com/msaldanha/realChain/block"
	"github.com/msaldanha/realChain/blockstore"
	"fmt"
	"os"
	"path/filepath"
	"github.com/spf13/viper"
	log "github.com/sirupsen/logrus"
)

const (
	cfgDataFolder = "ledger.datafolder"
	cfgChainFile = "ledger.chain"
	cfgAccountsFile = "ledger.accounts"
)

type Node struct {
	ld *ledger.Ledger
}

func New() (*Node) {
	return &Node{}
}

func (n *Node) Run(config *viper.Viper) {

	ld, err := createLedger(config)
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

func createLedger(config *viper.Viper) (*ledger.Ledger, error) {

	bklStoreOptions := prepareOptions("BlockChain",
		filepath.Join(config.GetString(cfgDataFolder), config.GetString(cfgChainFile)))
	blkStore := keyvaluestore.NewBoltKeyValueStore()
	err := blkStore.Init(bklStoreOptions)
	if err != nil {
		log.Fatal("Failed to init ledger chain: " + err.Error())
	}

	accStoreOptions := prepareOptions("Accounts",
		filepath.Join(config.GetString(cfgDataFolder), config.GetString(cfgAccountsFile)))
	as := keyvaluestore.NewBoltKeyValueStore()
	err = as.Init(accStoreOptions)
	if err != nil {
		log.Fatal("Failed to init ledger accounts" + err.Error())
	}

	val := block.NewBlockValidatorCreator()
	bs := blockstore.New(blkStore, val)
	ld := ledger.New()
	ld.Use(bs, as)
	if bs.IsEmpty() {
		if _, err := ld.Initialize(10000); err != nil {
			return nil, err
		}
	}
	return ld, nil
}

func prepareOptions(bucketName, filepath string) *keyvaluestore.BoltKeyValueStoreOptions {
	options := &keyvaluestore.BoltKeyValueStoreOptions{DbFile: filepath, BucketName: bucketName}
	return options
}