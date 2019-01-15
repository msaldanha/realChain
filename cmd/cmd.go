package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"github.com/msaldanha/realChain/config"
)

var rootCmd *cobra.Command
var cfg *viper.Viper

func init() {
	cfg = viper.New()

	cfg.SetConfigType("yaml")
	cfg.SetConfigName("config")
	cfg.AddConfigPath("/etc/realchain/")
	cfg.AddConfigPath("$HOME/.realchain")
	cfg.AddConfigPath(".")

	cfg.SetDefault(config.CfgDataFolder, "./")
	cfg.SetDefault(config.CfgLedgerChainFile, "chain.db")
	cfg.SetDefault(config.CfgLedgerAddressesFile, "addresses.db")
	cfg.SetDefault(config.CfgWalletChainFile, "wchain.db")
	cfg.SetDefault(config.CfgWalletAddressesFile, "waddresses.db")
	cfg.SetDefault(config.CfgNodeServer, "localhost:1300")
	cfg.SetDefault(config.CfgUdpServer, "localhost:1200")

	err := cfg.ReadInConfig()
	if err != nil {
		log.Fatalf("Error loading config file: %s \n", err)
	}

	rootCmd = &cobra.Command{Use: "realChain"}
	rootCmd.AddCommand(versionCmd)

	nodeCmd.AddCommand(nodeServerCmd)
	rootCmd.AddCommand(nodeCmd)

	ledgerCmd.AddCommand(ledgerInitCmd)
	rootCmd.AddCommand(ledgerCmd)


	walletCmd.AddCommand(walletListAddrsCmd)
	walletCmd.AddCommand(walletListAddressStatementCmd)
	walletCmd.AddCommand(walletSendCmd)
	walletCmd.AddCommand(walletCreateAddressCmd)
	rootCmd.AddCommand(walletCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of realChain",
	Long:  `Print the version number of realChain`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("realChain v0.1 -- HEAD")
	},
}

func New() *cobra.Command {
	return rootCmd
}