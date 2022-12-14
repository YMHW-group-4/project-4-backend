package commands

import (
	"backend/blockchain"
	"time"
)

func CreateBlockchain() blockchain.Blockchain{
	transaction := blockchain.Transaction{
		PubKeySender: "First Sender",
		PubKeyReceiver: "First Receiver",
		Amount: 00.00,
		Id: "First ID",
		Timestamp: time.Now(),
	}
	var transactions []blockchain.Transaction
	transactions = append(transactions, transaction)
	return blockchain.CreateBlockchain(transactions)
}

func NewTransaction() {

}

func ShowTransaction(){

}

func AddBlock(){

}

func CheckBlock(){

}

func ShowAllBlocks() {
	blockchain.ReadBlocks()
}
