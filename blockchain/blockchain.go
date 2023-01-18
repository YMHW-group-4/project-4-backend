package blockchain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// dumpFile the file name to whom the Blockchain should be written.
const dumpFile string = "blockchain.json"

// Blockchain holds all the blocks in the Blockchain.
type Blockchain struct {
	Blocks []Block `json:"blocks"`
	mp     *mempool
	am     *accountModel
}

// NewBlockchain creates a new Blockchain.
func NewBlockchain() *Blockchain {
	return &Blockchain{
		Blocks: make([]Block, 0),
		mp:     newMempool(),
		am:     newAccountModel(),
	}
}

// Init initializes the account model of the Blockchain.
func (b *Blockchain) Init() {
	if len(b.Blocks) > 0 {
		b.am.fromBlocks(b.Blocks...)
	}
}

func (b *Blockchain) AddBlock(block Block) error {
	if err := b.validate(block); err != nil {
		return err
	}

	b.Blocks = append(b.Blocks, block)

	return nil
}

func (b *Blockchain) CreateTransaction(sender string, receiver string, amount float32) error {
	if !b.am.exists(sender) {

	}

	tx, err := b.am.get(sender)
	if err != nil {
		return err
	}

	_ = Transaction{
		Sender:    sender,
		Receiver:  receiver,
		Signature: "", // TODO
		Amount:    amount,
		Nonce:     tx.transactions,
		Timestamp: time.Now().Unix(),
	}

	return nil
}

// TODO
func (b *Blockchain) validate(block Block) error {
	last := b.Blocks[len(b.Blocks)-1]

	if res := bytes.Compare(last.hash(), block.PrevHash); res != 0 {
		return fmt.Errorf("%w, %s", errInvalidBlock, "hash does not match")
	}

	if last.Timestamp > block.Timestamp {
		return fmt.Errorf("%w, %s", errInvalidBlock, "invalid timestamp")
	}

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
