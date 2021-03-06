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
	Use:   "init [amount] [address file]",
	Short: "Initializes a new ledger (transaction chain) with [amount] and put genesis address into [address file]",
	Long:  `Initializes a new ledger (transaction chain) with [amount] and put genesis address into [address file], creating a new genesis transaction`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Printf("Expected [amount] [address file]\n")
			os.Exit(1)
			return
		}

		startAmount, err := strconv.ParseFloat(args[0], 64)
		if err != nil {
			fmt.Printf("Failed to initialize the Ledger: %s\n", err)
			os.Exit(1)
		}
		if startAmount == 0 {
			fmt.Println("Failed to initialize the Ledger: amount must be > 0")
			os.Exit(1)
		}

		addressFile := args[1]

		bklStoreOptions := &keyvaluestore.BoltKeyValueStoreOptions{
			DbFile: filepath.Join(cfg.GetString(config.CfgDataFolder), cfg.GetString(config.CfgLedgerChainFile)),
			BucketName: "TxChain",
		}
		txStore := keyvaluestore.NewBoltKeyValueStore()
		err = txStore.Init(bklStoreOptions)
		if err != nil {
			log.Fatal("Failed to init ledger chain: " + err.Error())
		}

		asOpts := &keyvaluestore.BoltKeyValueStoreOptions{DbFile: filepath.Join(cfg.GetString(config.CfgDataFolder),
			addressFile), BucketName: "Addresses"}

		as := keyvaluestore.NewBoltKeyValueStore()

		err = as.Init(asOpts)
		if err != nil {
			fmt.Printf("Address store init failed: %s \n", err)
			os.Exit(1)
		}

		val := ledger.NewValidatorCreator()
		ts := ledger.NewTransactionStore(txStore, val)
		ld := ledger.NewLocalLedger(ts)
		if !ts.IsEmpty() {
			fmt.Println("Ledger already initialized")
			os.Exit(1)
		}
		if len(args) == 0 {
			fmt.Printf("Failed to initialize the Ledger: expected initial amount\n")
			os.Exit(1)
		}

		tx, addr, err := ledger.CreateGenesisTransaction(startAmount)
		if  err != nil {
			fmt.Printf("Failed to initialize the Ledger: %s\n", err)
			os.Exit(1)
		}

		err = as.Put(addr.Address, addr.ToBytes())
		if  err != nil {
			fmt.Printf("Failed to save genesis address: %s\n", err)
			os.Exit(1)
		}

		err = ld.Initialize(tx)
		if  err != nil {
			fmt.Printf("Failed to initialize the Ledger: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("Ledger successfuly initialized. Genesis address: %s, Start balance: %f\n", addr.Address, startAmount)
	},
}