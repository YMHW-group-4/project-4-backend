package blockchain

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func CreateBlockchain(transactions []Transaction) Blockchain{
	GenesisBlock := CreateGenesisBlock(transactions)
	var blocks []Block
	blocks = append(blocks, GenesisBlock)
	blockchain := Blockchain{
		Blocks: blocks,
	}
	file,_  :=json.MarshalIndent(blockchain, "", " ")
	_ = ioutil.WriteFile("../data/blockchain.json", file, 0644)
	return blockchain
}

func ReadBlocks() {
	dat, err := os.ReadFile("../data/blockchain.json")
	check(err)
	var res map[string]interface{}
	json.Unmarshal([]byte(dat), &res)
	blocks,_ := json.MarshalIndent(res["Blocks"], "", " ")
	fmt.Printf("The blockchain contains the following blocks: %s\n\n", string(blocks))
}

func addBlock(transactions []Transaction) {

}

func readBlockChain() {

}

func (blockchain *Blockchain) AddBlock(block Block) {
	blockchain.Blocks = append(blockchain.Blocks, block)
}
