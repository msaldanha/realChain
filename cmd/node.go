package cmd

import (
	"github.com/spf13/cobra"
	"github.com/msaldanha/realChain/node"
)

var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Node related commands",
	Long:  `Node related commands`,
}

var nodeServerCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start serving. Starts rest and udp servers",
	Long:  `Start serving. Starts rest and udp servers`,
	Run: func(cmd *cobra.Command, args []string) {
		node := node.New(cfg)
		node.Run()
	},
}

var nodeInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Init node configuration.",
	Long:  `Init node configuration creating node address keys.`,
	Run: func(cmd *cobra.Command, args []string) {
		node := node.New(cfg)
		node.Init()
	},
}
