package blockchain

import (
	"encoding/json"
	"os"
)

// dumpFile the file name to whom the Blockchain should be written.
const dumpFile string = "blockchain.json"

// Blockchain holds all the blocks in the Blockchain.
type Blockchain struct {
	Blocks []Block `json:"blocks"`
}

// NewBlockchain creates a new Blockchain.
func NewBlockchain() *Blockchain {
	return &Blockchain{
		Blocks: make([]Block, 0),
	}
}

func (b *Blockchain) AddBlock(block Block) error {
	if err := b.verify(block); err != nil {
		return err
	}

	b.Blocks = append(b.Blocks, block)

	return nil
}

func (b *Blockchain) verify(block Block) error {
	// TODO

	return nil
}

func (b *Blockchain) CreateGenesis() {
	// FIXME
	block, _ := CreateBlock([]Transaction{}, []byte("test"))
	b.Blocks = append(b.Blocks, block)
}

// FromFile returns all blocks that are written to the dumpfile.
func (b *Blockchain) FromFile() ([]Block, error) {
	var blockchain Blockchain

	data, err := os.ReadFile(dumpFile)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &blockchain)
	if err != nil {
		return nil, err
	}

	return blockchain.Blocks, nil
}

// DumpJSON writes the current Blockchain to a JSON file.
func (b *Blockchain) DumpJSON() error {
	data, err := json.MarshalIndent(b, "", " ")
	if err != nil {
		return err
	}

	if err = os.WriteFile(dumpFile, data, 0600); err != nil { //nolint
		return err
	}

	return nil
}
