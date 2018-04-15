package cmd

import (
	"github.com/spf13/cobra"
	"gopkg.in/resty.v1"
	"github.com/msaldanha/realChain/config"
	"fmt"
	"os"
	"bytes"
	"encoding/json"
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
		resp, err := resty.R().Get(getApiUrl("/wallet/addresses"))
		if err != nil {
			fmt.Printf("List addresses failed: %s (%s)", err, string(resp.Body()))
			os.Exit(1)
			return
		}
		fmt.Printf("Addresses: \n%s\n", getPrettyJson(resp.Body()))
	},
}

var walletCreateAddressCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates an address",
	Long:  `Creates an address`,
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := resty.R().Post(getApiUrl("/wallet/addresses"))
		if err != nil {
			fmt.Printf("Create address failed: %s (%s)", err, string(resp.Body()))
			os.Exit(1)
			return
		}
		fmt.Printf("Address: \n%s\n", getPrettyJson(resp.Body()))
	},
}

var walletListAddressStatementCmd = &cobra.Command{
	Use:   "statement [account address]",
	Short: "Lists all transactions for [address]",
	Long:  `Lists all transactions for [address]`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Printf("Expected [address]\n")
			os.Exit(1)
			return
		}
		resp, err := resty.R().Get(getApiUrl("/wallet/addresses" + "/" + args[0] + "/statement"))
		if err != nil {
			fmt.Printf("List address statement failed: %s (%s)", err, string(resp.Body()))
			os.Exit(1)
			return
		}
		fmt.Printf("Statement for address %s : \n%s\n", args[0], getPrettyJson(resp.Body()))
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
		body := fmt.Sprintf(`{"from": "%s","to": "%s","amount": %f}`, args[0], args[1], amount)
		resp, err := resty.R().
			SetBody(body).
			Post(getApiUrl("/wallet/tx"))
		if err != nil {
			fmt.Printf("Send transaction failed: %s (%s)", err, string(resp.Body()))
			os.Exit(1)
			return
		}
		fmt.Printf("Send transaction created : \n%s\n", getPrettyJson(resp.Body()))
	},
}

func getPrettyJson(jsonBytes []byte) string {
	var prettyJSON bytes.Buffer
	json.Indent(&prettyJSON, jsonBytes, "", "\t")
	return prettyJSON.String()
}

func getApiUrl(resource string) string {
	api := cfg.GetString(config.CfgRestServer)
	return "http://" + api + resource
}
