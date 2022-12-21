package cli

import (
	"backend/commands"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
)

var createBlockchainCmd = &cobra.Command{
	Use:   "createblockchain",
	Short: "Creates a first blockchain",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		res := commands.CreateBlockchain()
		resJson, _ := json.MarshalIndent(res, "", "   ")
		fmt.Println(string(resJson))
	},
}

var readBlockchainCmd = &cobra.Command{
	Use:   "readblockchain",
	Short: "Shows the complete blockchain present on system(node)",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		commands.ReadBlockchain()
	},
}

var readBlocksCmd = &cobra.Command{
	Use:   "readblocks",
	Short: "Shows all the blocks in the blockchain",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		commands.ShowAllBlocks()
	},
}

var readTransactions = &cobra.Command{
	Use:   "readtransactions",
	Short: "Shows all the transactions in the latest block",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		commands.ReadTransactions()
	},
}

func init() {
	rootCmd.AddCommand(createBlockchainCmd)
	rootCmd.AddCommand(readBlocksCmd)
	rootCmd.AddCommand(readBlockchainCmd)
	rootCmd.AddCommand(readTransactions)
}
