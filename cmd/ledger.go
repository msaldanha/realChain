package cmd

import (
	"github.com/spf13/cobra"
	"path/filepath"
	"github.com/msaldanha/realChain/keyvaluestore"
	"log"
	"github.com/msaldanha/realChain/transaction"
	"github.com/msaldanha/realChain/transactionstore"
	"github.com/msaldanha/realChain/ledger"
	"fmt"
	"os"
	"strconv"
	"github.com/msaldanha/realChain/config"
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
			DbFile: filepath.Join(cfg.GetString(config.CfgDataFolder), cfg.GetString(config.CfgChainFile)),
			BucketName: "TxChain",
		}
		txStore := keyvaluestore.NewBoltKeyValueStore()
		err := txStore.Init(bklStoreOptions)
		if err != nil {
			log.Fatal("Failed to init ledger chain: " + err.Error())
		}

		accStoreOptions := &keyvaluestore.BoltKeyValueStoreOptions{
			DbFile: filepath.Join(cfg.GetString(config.CfgDataFolder), cfg.GetString(config.CfgAccountsFile)),
			BucketName: "Accounts",
		}

		as := keyvaluestore.NewBoltKeyValueStore()
		err = as.Init(accStoreOptions)
		if err != nil {
			log.Fatal("Failed to init ledger accounts" + err.Error())
		}

		val := transaction.NewValidatorCreator()
		bs := transactionstore.New(txStore, val)
		ld := ledger.New()
		ld.Use(bs, as)
		if !bs.IsEmpty() {
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
		tx, err := ld.Initialize(startAmount)
		if  err != nil {
			fmt.Printf("Falied to initialize the Ledger: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("Ledger successfuly initialized. Genesis account: %s, Start balance: %f", string(tx.Account), startAmount)
	},
}