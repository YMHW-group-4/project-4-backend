package blockchain

import (
	"encoding/json"
	"io/ioutil"
)

type Blockchain struct {
	Blocks []Block
}

func CreateBlockchain(transactions []Transaction) Blockchain{
	GenesisBlock := CreateGenesisBlock(transactions)
	var blocks []Block
	blocks = append(blocks, GenesisBlock)
	blockchain := Blockchain{
		Blocks: blocks,
	}
	file,_  :=json.MarshalIndent(blockchain, "", " ")
	_ = ioutil.WriteFile("../data/blockchain", file, 0644)
	return blockchain
}

func (blockchain *Blockchain) AddBlock(block Block) {
	blockchain.Blocks = append(blockchain.Blocks, block)
}
