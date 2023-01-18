package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"time"
)

// errInvalidBlock is the base error when a block is invalid.
var errInvalidBlock = errors.New("invalid block")

// Block represents a singular block of the blockchain.
type Block struct {
	Validator    string        `json:"validator"`
	MerkleRoot   []byte        `json:"merkleRoot"`
	PrevHash     []byte        `json:"prevHash"`
	Height       int           `json:"height"`
	Timestamp    int64         `json:"timestamp"`
	Transactions []Transaction `json:"transactions"`
}

func CreateBlock(transactions []Transaction, prevHash []byte) (Block, error) {
	t, err := newMerkleTree(hashTransactions(transactions))
	if err != nil {
		return Block{}, err
	}

	return Block{
		MerkleRoot:   t.root.hash,
		PrevHash:     prevHash,
		Height:       len(transactions),
		Timestamp:    time.Now().Unix(),
		Transactions: transactions,
	}, nil
}

// string returns the block as a string.
func (b Block) string() string {
	return fmt.Sprintf("%v", b)
}

// hash returns the hash of the block.
func (b Block) hash() []byte {
	h := sha256.New()
	h.Write([]byte(b.string()))

	return h.Sum(nil)
}
