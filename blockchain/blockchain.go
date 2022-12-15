package blockchain

import (
	"encoding/json"
	"fmt"
	"os"
)

type Blockchain struct {
	Blocks []Block
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func CreateBlockchain(transactions []Transaction) Blockchain {
	GenesisBlock := CreateGenesisBlock(transactions)
	var blocks []Block
	blocks = append(blocks, GenesisBlock)
	blockchain := Blockchain{
		Blocks: blocks,
	}
	writeToFile(blockchain)
	return blockchain
}

func ReadBlocks() {
	dat := getBlockchainFromFile()
	blocks, _ := json.MarshalIndent(dat.Blocks, "", " ")
	fmt.Printf("The blockchain contains the following blocks: %s\n\n", string(blocks))
}

func AddBlocks(block Block) {
	dat := getBlockchainFromFile()
	dat.Blocks = append(dat.Blocks, block)
	writeToFile(dat)
}

func AddTransAction(transaction Transaction) {
	dat := getBlockchainFromFile()
	latestBlock := len(dat.Blocks) - 1
	dat.Blocks[latestBlock].Transactions = append(dat.Blocks[latestBlock].Transactions, transaction)
	writeToFile(dat)
}

func ReadBlockChain() {
	dat := getBlockchainFromFile()
	blockchain, _ := json.MarshalIndent(dat, "", " ")
	fmt.Printf("The blockchain on file: %s\n", string(blockchain))
}

func writeToFile(blockchain Blockchain) {
	file, _ := json.MarshalIndent(blockchain, "", " ")
	_ = os.WriteFile("../data/blockchain.json", file, 0o644)
}

func getBlockchainFromFile() Blockchain {
	dat, err := os.ReadFile("../data/blockchain.json")
	check(err)
	var blockchain Blockchain
	json.Unmarshal(dat, &blockchain)
	return blockchain
}
