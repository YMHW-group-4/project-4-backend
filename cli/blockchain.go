package cli

import (
	"backend/commands"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
)

var createBlockchainCmd = &cobra.Command{
	Use: "createblockchain",
	Short: "Creates a first blockchain",
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		commands.CreateBlockchain()
		res := commands.CreateBlockchain()
		resJson,_ := json.MarshalIndent(res, "","   ")
		fmt.Println(string(resJson))
	},
}

var readBlocksCmd = &cobra.Command {
	Use: "readblocks",
	Short: "Shows all the blocks in the blockchain",
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		commands.ShowAllBlocks()
	},
}

func init() {
	rootCmd.AddCommand(createBlockchainCmd)
	rootCmd.AddCommand(readBlocksCmd)
}