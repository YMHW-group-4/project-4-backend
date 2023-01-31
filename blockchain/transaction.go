package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"
)

// errInvalidTransaction is the base error when a transaction is invalid.
var errInvalidTransaction = errors.New("invalid transaction")

// Transaction represents a transaction within the blockchain.
type Transaction struct {
	Sender    string  `json:"sender"`
	Receiver  string  `json:"receiver"`
	Signature string  `json:"signature"`
	Amount    float64 `json:"amount"`
	Nonce     uint64  `json:"nonce"`
	Timestamp int64   `json:"timestamp"`
}

// string returns the transaction as a string.
func (t Transaction) string() string {
	return fmt.Sprintf("%v", t)
}

// hash returns the hash of the transaction.
func (t Transaction) hash() []byte {
	h := sha256.New()
	h.Write([]byte(t.string()))

	return h.Sum(nil)
}

// hashTransactions returns the hash of all given transactions.
func hashTransactions(transactions []Transaction) [][]byte {
	data := make([][]byte, 0, len(transactions))

	for _, t := range transactions {
		data = append(data, t.hash())
	}

	return data
}
