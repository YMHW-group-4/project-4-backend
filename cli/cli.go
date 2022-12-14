package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command {
	Use: "Hoin",
	Short: "Hoin - CLI tool for Honey blockchain",
	Long: `Hoin CLI is the Command line interface tool for the Hanzemoney (Honey) Blockchain.
	One can use Hoin CLI to interact with the blockchain, mainly around transaction, memorypool and the accountmodel`,
	Run: func(cmd *cobra.Command, args []string){

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "An error occured: '%s'", err)
		os.Exit(1)
	}
}
