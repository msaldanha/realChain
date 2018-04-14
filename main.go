package main

import (
	"github.com/msaldanha/realChain/node"
	"github.com/spf13/viper"
	"fmt"
)

func main() {

	config := viper.New()

	config.SetConfigType("yaml")
	config.SetConfigName("config")
	config.AddConfigPath("/etc/realchain/")
	config.AddConfigPath("$HOME/.realchain")
	config.AddConfigPath(".")

	config.SetDefault("ledger.datafolder", "./")
	config.SetDefault("ledger.chain", "chain.db")
	config.SetDefault("ledger.accounts", "accounts.db")

	err := config.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	node := node.New()
	node.Run(config)
}