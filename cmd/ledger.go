package cmd

import (
	"fmt"
	"github.com/msaldanha/realChain/config"
	"github.com/msaldanha/realChain/keyvaluestore"
	"github.com/msaldanha/realChain/ledger"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

var ledgerCmd = &cobra.Command{
	Use:   "ledger",
	Short: "Ledger related commands",
	Long:  `Ledger related commands`,
}

var ledgerInitCmd = &cobra.Command{
	Use:   "init [amount]",
	Short: "Initializes a new ledger (transaction chain) with [amount]",
	Long:  `Initializes a new ledger (transaction chain) with [amount], creating a new genesis transaction`,
	Run: func(cmd *cobra.Command, args []string) {

		bklStoreOptions := &keyvaluestore.BoltKeyValueStoreOptions{
			DbFile: filepath.Join(cfg.GetString(config.CfgDataFolder), cfg.GetString(config.CfgLedgerChainFile)),
			BucketName: "TxChain",
		}
		txStore := keyvaluestore.NewBoltKeyValueStore()
		err := txStore.Init(bklStoreOptions)
		if err != nil {
			log.Fatal("Failed to init ledger chain: " + err.Error())
		}

		accStoreOptions := &keyvaluestore.BoltKeyValueStoreOptions{
			DbFile: filepath.Join(cfg.GetString(config.CfgDataFolder), cfg.GetString(config.CfgLedgerAddressesFile)),
			BucketName: "Addresses",
		}

		as := keyvaluestore.NewBoltKeyValueStore()
		err = as.Init(accStoreOptions)
		if err != nil {
			log.Fatal("Failed to init ledger addresses" + err.Error())
		}

		val := ledger.NewValidatorCreator()
		ts := ledger.NewTransactionStore(txStore, val)
		//ld := ledger.NewLocalLedger(ts)
		if !ts.IsEmpty() {
			fmt.Println("Ledger already initialized")
			os.Exit(1)
		}
		if len(args) == 0 {
			fmt.Printf("Falied to initialize the Ledger: expected initial amount\n")
			os.Exit(1)
		}
		startAmount, err := strconv.ParseFloat(args[0], 64)
		if err != nil {
			fmt.Printf("Falied to initialize the Ledger: %s\n", err)
			os.Exit(1)
		}
		if startAmount == 0 {
			fmt.Println("Falied to initialize the Ledger: amount must be > 0")
			os.Exit(1)
		}
		//TODO: fix me
		//err := ld.Initialize(startAmount)
		//if  err != nil {
		//	fmt.Printf("Falied to initialize the Ledger: %s\n", err)
		//	os.Exit(1)
		//}
		//fmt.Printf("Ledger successfuly initialized. Genesis address: %s, Start balance: %f", string(tx.Address), startAmount)
	},
}