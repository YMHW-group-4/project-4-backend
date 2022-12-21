package blockchain

import (
	"encoding/json"
	"fmt"
	"os"
)

type Blockchain struct {
	Blocks []Block
}

// check is a standard error checking method.
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// CreateBlockchain creates a blockchain and writes to JSON file.
// The blockchain contains a Genesis block with genesis hash and dummy transaction.
func CreateBlockchain(transactions []Transaction) Blockchain {
	GenesisBlock := CreateGenesisBlock(transactions)
	var blocks []Block
	blocks = append(blocks, GenesisBlock)
	blockchain := Blockchain{
		Blocks: blocks,
	}
	return blockchain
}

// ReadBlockChain gets the whole blockchain from JSON file and prints it whole.
func ReadBlockChain() {
	dat := getBlockchainFromFile()
	blockchain, _ := json.MarshalIndent(dat, "", " ")
	fmt.Printf("The blockchain on file: %s\n", string(blockchain))
}

// AddBlockToBlockchain Adds block to the blockchain.
// Block is first made and then given as parameter. Also writes to JSON blockchain file.
func (blockchain *Blockchain) AddBlockToBlockchain(transactions []Transaction) {
	block := CreateBlock("Test")
	block.Transactions = transactions
	blockchain.Blocks = append(blockchain.Blocks, block)
	blockchain.WriteToFile()
}

// ReadBlocks shows all the blocks on the blockchain, with all the transactions in them.
func ReadBlocks() {
	dat := getBlockchainFromFile()
	blocks, _ := json.MarshalIndent(dat.Blocks, "", " ")
	fmt.Printf("The blockchain contains the following blocks: %s\n\n", string(blocks))
}

// AddTransAction adds a transaction to the blockchain.
// A transaction has to be made first and then given in the method as parameter. Also writes to JSON blockchain file.
//func (blockchain *Blockchain) AddTransAction(transaction Transaction) {
//	blockchain.Transactions = append(blockchain.Transactions, transaction)
//}

// ReadTransactions gets and shows all the transactions present on the latest block.
func ReadTransactions() []Transaction {
	dat := getBlockchainFromFile()
	latestBlock := len(dat.Blocks) - 1
	return dat.Blocks[latestBlock].Transactions
}

// WriteToFile writes the blockchain to the JSON blockchain file.
// This method is called upon everytime a block or transaction is added to the blockchain.
func (blockchain *Blockchain) WriteToFile() {
	file, _ := json.MarshalIndent(blockchain, "", " ")
	_ = os.WriteFile("../data/blockchain.json", file, 0o644)
}

// getBlockchainFromFile gets the whole blockchain from the JSON blockchain file.
// This method is called upon everytime something is read from the blockchain.
func getBlockchainFromFile() Blockchain {
	dat, err := os.ReadFile("../data/blockchain.json")
	check(err)
	var blockchain Blockchain
	json.Unmarshal(dat, &blockchain)
	return blockchain
}
