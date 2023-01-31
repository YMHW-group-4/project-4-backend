package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// errInvalidBlock is the base error when a block is invalid.
var errInvalidBlock = errors.New("invalid block")

// Block represents a singular block of the blockchain.
type Block struct {
	Validator    string        `json:"validator"`
	MerkleRoot   string        `json:"merkleRoot"`
	PrevHash     string        `json:"prevHash"`
	Height       uint64        `json:"height"`
	Timestamp    int64         `json:"timestamp"`
	Transactions []Transaction `json:"transactions"`
}

// createBlock creates a new Block.
func createBlock(validator string, prevHash []byte, transactions []Transaction) (Block, error) {
	if len(transactions) == 0 {
		return Block{}, fmt.Errorf("%w: zero transactions", errInvalidBlock)
	}

	t, err := newMerkleTree(hashTransactions(transactions))
	if err != nil {
		return Block{}, err
	}

	return Block{
		Validator:    validator,
		MerkleRoot:   hex.EncodeToString(t.root.hash),
		PrevHash:     hex.EncodeToString(prevHash),
		Height:       uint64(len(transactions)),
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
