package node

import (
	"fmt"
	"github.com/msaldanha/realChain/address"
	"github.com/msaldanha/realChain/errors"
	"github.com/msaldanha/realChain/config"
	"github.com/msaldanha/realChain/consensus"
	"github.com/msaldanha/realChain/keyvaluestore"
	"github.com/msaldanha/realChain/ledger"
	"github.com/msaldanha/realChain/peerdiscovery"
	"github.com/msaldanha/realChain/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net"
	"os"
	"path/filepath"
)

const ErrLedgerNotInitialized = errors.Error("ledger not initialized")

type Node struct {
	addrDb keyvaluestore.Storer
	ld     ledger.Ledger
	cfg    *viper.Viper
}

func New(cfg *viper.Viper) *Node {
	return &Node{cfg: cfg}
}

func (n *Node) Run() {
	log.Info("Loading address db.")
	err := n.loadAddrDb()
	checkError(err)

	log.Info("Creating ledger.")
	err = n.createLedger()
	checkError(err)

	log.Info("Creating server.")
	srv, err := n.createServer()
	checkError(err)

	serverCh := make(chan error)

	go func() { serverCh <- srv.Run() }()

	log.Info("Ready.")
	er := <-serverCh
	checkError(er)

	fmt.Println("Done.")
}

func (n *Node) Init() {
	log.Info("Loading address db.")
	err := n.loadAddrDb()
	checkError(err)

	log.Info("Creating address.")
	addr, err := address.NewAddressWithKeys()
	checkError(err)

	log.Info("Saving address.")
	err = n.addrDb.Put(addr.Address, addr.ToBytes())
	checkError(err)

	fmt.Printf("Done. Node address: %s\n", addr.Address)
}

func (n *Node) getNodeAddr() *address.Address {
	addrs, err := n.addrDb.GetAll()
	checkError(err)
	addr := address.NewAddressFromBytes(addrs[0])
	return addr
}

func (n *Node) loadAddrDb() error {
	addrOptions := prepareOptions(config.AddressBucket,
		filepath.Join(n.cfg.GetString(config.CfgDataFolder), n.cfg.GetString(config.CfgNodeAddressesFile)))

	addrDb := keyvaluestore.NewBoltKeyValueStore()

	err := addrDb.Init(addrOptions)
	if err != nil {
		return err
	}

	n.addrDb = addrDb
	return nil
}

func checkError(err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}

func (n *Node) createLedger() error {
	txOptions := prepareOptions(config.TxBucket,
		filepath.Join(n.cfg.GetString(config.CfgDataFolder), n.cfg.GetString(config.CfgLedgerChainFile)))
	txDb := keyvaluestore.NewBoltKeyValueStore()
	err := txDb.Init(txOptions)
	if err != nil {
		log.Fatal("Failed to init ledger chain: " + err.Error())
	}

	val := ledger.NewValidatorCreator()
	bs := ledger.NewTransactionStore(txDb, val)
	if bs.IsEmpty() {
		return ErrLedgerNotInitialized
	}

	n.ld = ledger.NewLocalLedger(bs)

	return nil
}

func (n *Node) createServer() (*server.Server, error) {
	addr := n.getNodeAddr()
	con := consensus.NewConsensus(n.ld, addr)
	dis := peerdiscovery.NewStaticDiscoverer(n.cfg)

	listener, err := net.Listen("tcp", n.cfg.GetString(config.CfgNodeServer))
	if err != nil {
		return nil, err
	}
	return server.New(n.ld, con, dis, listener), nil
}

func prepareOptions(bucketName, filepath string) *keyvaluestore.BoltKeyValueStoreOptions {
	options := &keyvaluestore.BoltKeyValueStoreOptions{DbFile: filepath, BucketName: bucketName}
	return options
}
