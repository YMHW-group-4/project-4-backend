package blockchain

import (
	"crypto/sha256"
	"time"
)

type Block struct {
	Transactions []Transaction
	Hash         []byte
	Timestamp    time.Time
}

func CreateGenesisBlock(transactions []Transaction) Block {
	h := Block{
		Transactions: transactions,
		Hash:         createGenesisHash(),
		Timestamp:    time.Now(),
	}
	return h
}

func createGenesisHash() []byte {
	s := "genesishash"
	h := sha256.New()
	h.Write([]byte(s))
	hash := h.Sum(nil)
	return hash
}
