package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"backend/util"
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

// newBlock creates a new Block.
func newBlock(validator string, prevHash []byte, transactions []Transaction) (Block, error) {
	if len(transactions) == 0 {
		return Block{}, fmt.Errorf("%w: zero transactions", errInvalidBlock)
	}

	t, err := newMerkleTree(hashTransactions(transactions))
	if err != nil {
		return Block{}, err
	}

	return Block{
		Validator:    validator,
		MerkleRoot:   util.HexEncode(t.root.hash),
		PrevHash:     util.HexEncode(prevHash),
		Height:       uint64(len(transactions)),
		Timestamp:    time.Now().Unix(),
		Transactions: transactions,
	}, nil
}

// string returns the block as a string.
func (b Block) string() string {
	return fmt.Sprintf("%v", b)
}

// Hash returns the hash of the Block.
func (b Block) Hash() []byte {
	h := sha256.New()
	h.Write([]byte(b.string()))

	return h.Sum(nil)
}

// Validate validates a singular Block.
func (b Block) Validate(last Block, validator string) error {
	// compare hashes
	if util.HexEncode(last.Hash()) != b.PrevHash {
		return fmt.Errorf("%w, %s", errInvalidBlock, "hash does not match")
	}

	// check timstamp
	if last.Timestamp > b.Timestamp {
		return fmt.Errorf("%w, %s", errInvalidBlock, "invalid timestamp")
	}

	// compare validator
	if b.Validator != validator {
		return fmt.Errorf("%w, %s", errInvalidBlock, "invalid validator")
	}

	// create new tree
	tr, err := newMerkleTree(hashTransactions(b.Transactions))
	if err != nil {
		return fmt.Errorf("%w, %s", errInvalidBlock, "failed to create tree")
	}

	// compare merkle root
	if util.HexEncode(tr.root.hash) != b.MerkleRoot {
		return fmt.Errorf("%w, %s", errInvalidBlock, "merkle root does not match")
	}

	// compare height
	if uint64(len(b.Transactions)) != b.Height {
		return fmt.Errorf("%w, %s", errInvalidBlock, "height does not match")
	}

	return nil
}
