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
	"github.com/msaldanha/realChain/network"
	"github.com/msaldanha/realChain/wallet"
	"github.com/msaldanha/realChain/consensus"
)

const ErrLedgerNotInitialized = Error.Error("ledger not initialized")

type Node struct {
	ld  ledger.Ledger
	cfg *viper.Viper
}

func New(cfg *viper.Viper) (*Node) {
	return &Node{cfg: cfg}
}

func (n *Node) Run() {
	ld, err := n.createLedger()
	checkError(err)

	wa, err := n.createWallet(ld)
	checkError(err)

	walletRestServer, err := NewWalletRestServer(wa, n.cfg.GetString(config.CfgWalletRestServer))
	checkError(err)

	ledgerRestServer, err := NewLedgerRestServer(ld, n.cfg.GetString(config.CfgNodeRestServer))
	checkError(err)

	net, err := n.createNetwork()
	checkError(err)

	con := consensus.New(net, ld)

	udpch := make(chan error)
	walletRestCh := make(chan error)
	ledgerRestCh := make(chan error)

	go func() { udpch <- net.Run() }()
	go func() { walletRestCh <- walletRestServer.Run() }()
	go func() { ledgerRestCh <- ledgerRestServer.Run() }()
	go con.Run()

	for {
		select {
		case eu := <-udpch:
			checkError(eu)
		case er := <-walletRestCh:
			checkError(er)
		case er := <-ledgerRestCh:
			checkError(er)
		}
	}
	fmt.Println("Done.")
}

func (n *Node) createNetwork() (*network.Network, error) {
	localUrl := n.cfg.GetString(config.CfgNodeNetworkLocal)
	net, err := network.NewNetwork(localUrl)
	if err != nil {
		return nil, err
	}
	peers := n.cfg.GetStringSlice(config.CfgNodeNetworkPeers)
	net.UsePeers(peers...)
	return net, nil
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

	val := transaction.NewValidatorCreator()
	bs := transactionstore.New(txDb, val)
	if bs.IsEmpty() {
		return nil, ErrLedgerNotInitialized
	}

	ld := ledger.NewLocalLedger(bs, addrDb)

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

	val := transaction.NewValidatorCreator()
	ts := transactionstore.New(txDb, val)

	wa := wallet.New(ts, addrDb, ld)

	return wa, nil
}

func prepareOptions(bucketName, filepath string) *keyvaluestore.BoltKeyValueStoreOptions {
	options := &keyvaluestore.BoltKeyValueStoreOptions{DbFile: filepath, BucketName: bucketName}
	return options
}
