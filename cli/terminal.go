package cli

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "Hoin",
	Short: "Hoin - CLI tool for Honey blockchain",
	Long: `Hoin CLI is the Command line interface tool for the Hanzemoney (Honey) Blockchain.
	One can use Hoin CLI to interact with the blockchain, mainly around transactions, memorypool and the accountmodel`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error().Err(err).Msg("cli: failed to execute")
	}
}
