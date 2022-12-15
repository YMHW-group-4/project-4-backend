package cli

import (
	"backend/blockchain"
	"backend/commands"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
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

var addBlock = &cobra.Command{
	Use:   "addblock",
	Short: "Adds block to the blockchain",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		commands.AddBlock(args[0])
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

var addTransactionToBlock = &cobra.Command{
	Use:   "addtransaction",
	Short: "Adds transaction to the latest block",
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		var amount float32
		if s, err := strconv.ParseFloat(args[2], 32); err == nil {
			amount = float32(s)
		}
		transaction := blockchain.Transaction{
			PubKeySender:   args[0],
			PubKeyReceiver: args[1],
			Amount:         amount,
			Id:             args[3],
		}
		commands.NewTransaction(transaction)
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
	rootCmd.AddCommand(addTransactionToBlock)
	rootCmd.AddCommand(addBlock)
	rootCmd.AddCommand(readTransactions)
}
