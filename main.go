package main

import (
	"github.com/msaldanha/realChain/node"
	"github.com/spf13/viper"
	log "github.com/sirupsen/logrus"
)

func main() {

	log.Info("Application starting")
	config := viper.New()

	config.SetConfigType("yaml")
	config.SetConfigName("config")
	config.AddConfigPath("/etc/realchain/")
	config.AddConfigPath("$HOME/.realchain")
	config.AddConfigPath(".")

	config.SetDefault("ledger.datafolder", "./")
	config.SetDefault("ledger.chain", "chain.db")
	config.SetDefault("ledger.accounts", "accounts.db")
	config.SetDefault("node.restserver", "localhost:1300")
	config.SetDefault("node.udpserver", "localhost:1200")

	err := config.ReadInConfig()
	if err != nil {
		log.Fatalf("Error loading config file: %s \n", err)
	}

	node := node.New()
	node.Run(config)
}