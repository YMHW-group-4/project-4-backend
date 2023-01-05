package blockchain

import (
	"crypto/sha256"
	"fmt"
	"time"
)

type Block struct {
	MerkleRoot   []byte        `json:"merkleRoot"`
	PrevHash     []byte        `json:"prevHash"`
	Height       int           `json:"height"`
	Timestamp    int64         `json:"timestamp"`
	Transactions []Transaction `json:"transactions"`
}

func CreateBlock(transactions []Transaction, prevHash []byte) (Block, error) {
	// FIXME find other way of hashing transactions for merkle tree (?)
	// perhaps pass type T to tree and do internal hashing in the tree.
	data := make([][]byte, 0)
	h := sha256.New()

	for _, t := range transactions {
		h.Write([]byte(fmt.Sprintf("%v", t)))
		data = append(data, h.Sum(nil))
	}

	t, err := newMerkleTree(data)
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
