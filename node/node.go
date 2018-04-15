package node

import (
	"github.com/msaldanha/realChain/config"
	"github.com/msaldanha/realChain/ledger"
	"github.com/msaldanha/realChain/keyvaluestore"
	"github.com/msaldanha/realChain/transaction"
	"github.com/msaldanha/realChain/transactionstore"
	"fmt"
	"os"
	"path/filepath"
	"github.com/spf13/viper"
	log "github.com/sirupsen/logrus"
	"github.com/msaldanha/realChain/Error"
)

const ErrLedgerNotInitialized            = Error.Error("ledger not initialized")

type Node struct {
	ld *ledger.Ledger
}

func New() (*Node) {
	return &Node{}
}

func (n *Node) Run(cfg *viper.Viper) {

	ld, err := createLedger(cfg)
	checkError(err)

	udp := NewUdpServer(ld, cfg.GetString(config.CfgUdpServer))
	rest := NewRestServer(ld, cfg.GetString(config.CfgRestServer))

	udpch := make(chan error)
	restch := make(chan error)

	go func() { udpch <- udp.Run() }()
	go func() { restch <- rest.Run() }()

	for {
		select {
		case eu := <-udpch:
			checkError(eu)
		case er := <-restch:
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

func createLedger(cfg *viper.Viper) (*ledger.Ledger, error) {
	bklStoreOptions := prepareOptions("TxChain",
		filepath.Join(cfg.GetString(config.CfgDataFolder), cfg.GetString(config.CfgChainFile)))
	txStore := keyvaluestore.NewBoltKeyValueStore()
	err := txStore.Init(bklStoreOptions)
	if err != nil {
		log.Fatal("Failed to init ledger chain: " + err.Error())
	}

	accStoreOptions := prepareOptions("Accounts",
		filepath.Join(cfg.GetString(config.CfgDataFolder), cfg.GetString(config.CfgAccountsFile)))
	as := keyvaluestore.NewBoltKeyValueStore()
	err = as.Init(accStoreOptions)
	if err != nil {
		log.Fatal("Failed to init ledger accounts" + err.Error())
	}

	val := transaction.NewValidatorCreator()
	bs := transactionstore.New(txStore, val)
	if bs.IsEmpty() {
		return nil, ErrLedgerNotInitialized
	}

	ld := ledger.New()
	ld.Use(bs, as)

	return ld, nil
}

func prepareOptions(bucketName, filepath string) *keyvaluestore.BoltKeyValueStoreOptions {
	options := &keyvaluestore.BoltKeyValueStoreOptions{DbFile: filepath, BucketName: bucketName}
	return options
}
