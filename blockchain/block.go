package blockchain

import (
	"crypto/sha256"
	"time"
)

type Block struct {
	Hash         []byte
	Timestamp    time.Time
	Transactions []Transaction
}

func CreateGenesisBlock(transactions []Transaction) Block {
	b := Block{
		Hash:         createGenesisHash(),
		Timestamp:    time.Now(),
		Transactions: transactions,
	}
	return b
}

func createGenesisHash() []byte {
	s := "genesishash"
	h := sha256.New()
	h.Write([]byte(s))
	hash := h.Sum(nil)
	return hash
}

func createHash(key string) []byte {
	s := key
	h := sha256.New()
	h.Write([]byte(s))
	hash := h.Sum(nil)
	return hash
}

func CreateBlock(key string) Block {
	b := Block{
		Hash:         createHash(key),
		Timestamp:    time.Now(),
		Transactions: []Transaction{},
	}
	return b
}
