package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/msaldanha/realChain/config"
	"github.com/msaldanha/realChain/keyvaluestore"
	"github.com/msaldanha/realChain/ledger"
	"github.com/msaldanha/realChain/wallet"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"os"
	"path/filepath"
	"strconv"
)

var walletCmd = &cobra.Command{
	Use:   "wallet",
	Short: "Wallet related commands",
	Long:  `Wallet related commands`,
}

var walletListAddrsCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all managed addresses",
	Long:  `Lists all managed addresses`,
	Run: func(cmd *cobra.Command, args []string) {
		wa := getWallet()
		addrs, err := wa.GetAddresses()
		if err != nil {
			fmt.Printf("List addresses failed: %s \n", err)
			os.Exit(1)
			return
		}

		fmt.Printf("Addresses: \n%s\n", getPrettyJson(addrs))
	},
}

var walletCreateAddressCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates an address",
	Long:  `Creates an address`,
	Run: func(cmd *cobra.Command, args []string) {
		wa := getWallet()
		addr, err := wa.CreateAddress()
		if err != nil {
			fmt.Printf("Create address failed: %s \n", err)
			os.Exit(1)
			return
		}

		fmt.Printf("Address: \n%s\n", getPrettyJson(addr))
	},
}

var walletListAddressStatementCmd = &cobra.Command{
	Use:   "statement [address]",
	Short: "Lists all transactions for [address]",
	Long:  `Lists all transactions for [address]`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Printf("Expected [address]\n")
			os.Exit(1)
			return
		}

		wa := getWallet()
		st, err := wa.GetAddressStatement(args[0])
		if err != nil {
			fmt.Printf("List address statement failed: %s \n", err)
			os.Exit(1)
			return
		}

		fmt.Printf("Statement for address %s : \n%s\n", args[0], getPrettyJson(st))
	},
}

var walletSendCmd = &cobra.Command{
	Use:   "send [FROM address] [TO address] [amount]",
	Short: "Sends [amount] from [FROM address] to [TO address]",
	Long:  `Sends [amount] from [FROM address] to [TO address]`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 3 {
			fmt.Printf("Expected [FROM address] [TO address] [amount]\n")
			os.Exit(1)
			return
		}

		amount, err := strconv.ParseFloat(args[2], 64)
		if err != nil {
			fmt.Printf("Failed reading [amount]: %s\n", err)
			os.Exit(1)
		}

		wa := getWallet()
		tx, err := wa.Transfer(args[0], args[1], amount)
		if err != nil {
			fmt.Printf("Send transaction failed: %s \n", err)
			os.Exit(1)
			return
		}

		fmt.Printf("Send transaction created : \n%s\n", getPrettyJson(tx))
	},
}

func getPrettyJson(v interface{}) string {
	var prettyJSON bytes.Buffer
	jsonBytes, _ := json.Marshal(v)
	_ = json.Indent(&prettyJSON, jsonBytes, "", "\t")
	return prettyJSON.String()
}

func getWallet() *wallet.Wallet {
	options := &keyvaluestore.BoltKeyValueStoreOptions{DbFile: filepath.Join(cfg.GetString(config.CfgDataFolder),
		cfg.GetString(config.CfgWalletAddressesFile)), BucketName: "Addresses"}

	as := keyvaluestore.NewBoltKeyValueStore()

	err := as.Init(options)
	if err != nil {
		fmt.Printf("Wallet address store init failed: %s ", err)
		os.Exit(1)
	}

	conn, err := grpc.Dial(cfg.GetString(config.CfgNodeServer), grpc.WithInsecure())
	if err != nil {
		fmt.Printf("Wallet connection to ledger failed: %s ", err)
		os.Exit(1)
	}

	ld := ledger.NewLedgerClient(conn)

	return wallet.New(as, ld)
}
