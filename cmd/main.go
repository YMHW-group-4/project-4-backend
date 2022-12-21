package main

import (
	"backend/blockchain"
	"backend/cli"
	"bufio"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
	"strings"
	"time"
)

// setupLogger checks whether the Stdout is a cli or not
// if so it sets the global log's writer to a ConsoleWriter.
func setupLogger() {
	if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
		logOutput := zerolog.ConsoleWriter{Out: os.Stderr}
		log.Logger = log.Output(logOutput)

		return
	}

	log.Logger = log.Output(os.Stderr)
}

// setLogLevel sets the global log level to either debug or info level.
func setLogLevel(debug bool) {
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

func setupBlockchain() blockchain.Blockchain {
	var transactions []blockchain.Transaction
	blockchain := blockchain.CreateBlockchain(transactions)
	return blockchain
}

func main() {
	//startup := time.Now()
	//config := getConfigFromEnv()
	//setLogLevel(config.debug)
	//setupLogger()
	//
	//log.Info().
	//	Str("version", version).
	//	Int("port", config.port).
	//	Bool("debug", config.debug).
	//	Msg("node: startup")
	//
	//node, err := newNode(config)
	//if err != nil {
	//	log.Fatal().Err(err).Msg("node: failed to create node")
	//}
	//
	//node.run()
	//
	//log.Info().
	//	TimeDiff("startup", time.Now(), startup).
	//	Msg("node: running")

	var approvedTransactions []blockchain.Transaction
	currentBlockchain := setupBlockchain()
	currentBlockchain.WriteToFile()
	cli.Execute()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter text: ")
		scanner.Scan()
		text := scanner.Text()

		switch text {
		case "addtransaction":
			fmt.Print("Enter sender, receiver, amount and id: ")
			scanner.Scan()
			text = scanner.Text()
			textArray := strings.Fields(text)
			var amount float32
			if s, err := strconv.ParseFloat(textArray[2], 32); err == nil {
				amount = float32(s)
			}
			transaction := blockchain.Transaction{
				PubKeySender:   textArray[0],
				PubKeyReceiver: textArray[1],
				Amount:         amount,
				Id:             textArray[3],
				Timestamp:      time.Now(),
			}
			approvedTransactions = append(approvedTransactions, transaction)
			fmt.Print("Current approved transactions: \n", approvedTransactions, "\n")
		case "addblock":
			currentBlockchain.AddBlockToBlockchain(approvedTransactions)
			fmt.Print("The block was added to the blockchain, and looks like this: \n", currentBlockchain.Blocks[len(currentBlockchain.Blocks)-1], "\n")
		}
	}
	//node.handleSigterm()
}
