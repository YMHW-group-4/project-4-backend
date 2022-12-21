package commands

import (
	"encoding/json"
	"fmt"
	"time"

	"backend/blockchain"
)

func CreateBlockchain() blockchain.Blockchain {
	transaction := blockchain.Transaction{
		PubKeySender:   "First Sender",
		PubKeyReceiver: "First Receiver",
		Amount:         00.00,
		Id:             "First ID",
		Timestamp:      time.Now(),
	}
	var transactions []blockchain.Transaction
	transactions = append(transactions, transaction)
	return blockchain.CreateBlockchain(transactions)
}

func ReadTransactions() {
	transactions := blockchain.ReadTransactions()
	transJson, _ := json.MarshalIndent(transactions, "", " ")
	fmt.Printf("The latest block contains the following transactions:\n %s\n", transJson)
}

func ReadBlockchain() {
	blockchain.ReadBlockChain()
}

func ShowAllBlocks() {
	blockchain.ReadBlocks()
}
