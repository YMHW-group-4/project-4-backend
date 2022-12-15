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
	b := Block{
		Transactions: transactions,
		Hash:         createGenesisHash(),
		Timestamp:    time.Now(),
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

func createHash() []byte {
	s := "secondhash"
	h := sha256.New()
	h.Write([]byte(s))
	hash := h.Sum(nil)
	return hash
}

func CreateBlock() Block {
	b := Block{
		Transactions: []Transaction{},
		Hash:         createHash(),
		Timestamp:    time.Now(),
	}
	return b
}
