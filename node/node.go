package node

import (
	"fmt"
	"github.com/msaldanha/realChain/errors"
	"github.com/msaldanha/realChain/config"
	"github.com/msaldanha/realChain/consensus"
	"github.com/msaldanha/realChain/keyvaluestore"
	"github.com/msaldanha/realChain/ledger"
	"github.com/msaldanha/realChain/peerdiscovery"
	"github.com/msaldanha/realChain/server"
	"github.com/msaldanha/realChain/wallet"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net"
	"os"
	"path/filepath"
)

const ErrLedgerNotInitialized = errors.Error("ledger not initialized")

type Node struct {
	ld  ledger.Ledger
	cfg *viper.Viper
}

func New(cfg *viper.Viper) *Node {
	return &Node{cfg: cfg}
}

func (n *Node) Run() {
	ld, err := n.createLedger()
	checkError(err)

	wa, err := n.createWallet(ld)
	checkError(err)

	walletRestServer, err := NewWalletRestServer(wa, n.cfg.GetString(config.CfgWalletRestServer))
	checkError(err)

	server, err := n.createServer(ld)
	checkError(err)

	walletRestCh := make(chan error)
	serverCh := make(chan error)

	go func() { walletRestCh <- walletRestServer.Run() }()
	go func() { serverCh <- server.Run() }()

	for {
		select {
		case er := <-walletRestCh:
			checkError(er)
		case er := <-serverCh:
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

func (n *Node) createLedger() (ledger.Ledger, error) {
	txOptions := prepareOptions("TxChain",
		filepath.Join(n.cfg.GetString(config.CfgDataFolder), n.cfg.GetString(config.CfgLedgerChainFile)))
	txDb := keyvaluestore.NewBoltKeyValueStore()
	err := txDb.Init(txOptions)
	if err != nil {
		log.Fatal("Failed to init ledger chain: " + err.Error())
	}

	addrOptions := prepareOptions("Addresses",
		filepath.Join(n.cfg.GetString(config.CfgDataFolder), n.cfg.GetString(config.CfgLedgerAddressesFile)))
	addrDb := keyvaluestore.NewBoltKeyValueStore()
	err = addrDb.Init(addrOptions)
	if err != nil {
		log.Fatal("Failed to init ledger addresses" + err.Error())
	}

	val := ledger.NewValidatorCreator()
	bs := ledger.NewTransactionStore(txDb, val)
	if bs.IsEmpty() {
		return nil, ErrLedgerNotInitialized
	}

	ld := ledger.NewLocalLedger(bs)

	return ld, nil
}

func (n *Node) createWallet(ld ledger.Ledger) (*wallet.Wallet, error) {
	txOptions := prepareOptions("TxChain",
		filepath.Join(n.cfg.GetString(config.CfgDataFolder), n.cfg.GetString(config.CfgWalletChainFile)))
	txDb := keyvaluestore.NewBoltKeyValueStore()
	err := txDb.Init(txOptions)
	if err != nil {
		log.Fatal("Failed to init wallet chain: " + err.Error())
	}

	addrOptions := prepareOptions("Addresses",
		filepath.Join(n.cfg.GetString(config.CfgDataFolder), n.cfg.GetString(config.CfgWalletAddressesFile)))
	addrDb := keyvaluestore.NewBoltKeyValueStore()
	err = addrDb.Init(addrOptions)
	if err != nil {
		log.Fatal("Failed to init wallet addresses" + err.Error())
	}

	val := ledger.NewValidatorCreator()
	ts := ledger.NewTransactionStore(txDb, val)

	wa := wallet.New(ts, addrDb, ld)

	return wa, nil
}

func (n *Node) createServer(ld ledger.Ledger) (*server.Server, error) {
	consensus := consensus.NewConsensus(ld)
	discovery := peerdiscovery.NewStaticDiscoverer()

	listener, err := net.Listen("tcp", n.cfg.GetString(config.CfgNodeServer))
	if err != nil {
		return nil, err
	}
	return server.New(ld, consensus, discovery, listener), nil
}

func prepareOptions(bucketName, filepath string) *keyvaluestore.BoltKeyValueStoreOptions {
	options := &keyvaluestore.BoltKeyValueStoreOptions{DbFile: filepath, BucketName: bucketName}
	return options
}
