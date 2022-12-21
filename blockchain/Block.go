package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"math/rand"
	"time"
)

type Block struct {
	prevHash    string
	transaction Transaction
	timestamp   int64
	nonce       int64
}

func (block Block) Hash() string {
	str, _ := json.Marshal(block)
	h := sha256.New()
	h.Write(str)
	return string(h.Sum(nil))
}

func CreateBlock(prevHash string, transaction Transaction) Block {
	nonce := rand.NewSource(time.Now().UnixNano())
	return Block{
		prevHash:    prevHash,
		transaction: transaction,
		timestamp:   time.Now().UnixMilli(),
		nonce:       nonce.Int63(),
	}
}
