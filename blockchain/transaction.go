package blockchain

import (
	"backend/crypto"
	"backend/util"
	"crypto/sha256"
	"errors"
	"fmt"
)

// TxType is the type of the Transaction.
type TxType string

const (
	Stake    TxType = "stake"
	Regular  TxType = "regular"
	Reward   TxType = "reward"
	Fee      TxType = "fee"
	Penalty  TxType = "penalty"
	Exchange TxType = "exchange"
)

// ErrInvalidTransaction is the base error when a transaction is invalid.
var ErrInvalidTransaction = errors.New("invalid transaction")

// Transaction represents a transaction within the blockchain.
type Transaction struct {
	Sender    string  `json:"sender"`
	Receiver  string  `json:"receiver"`
	Signature string  `json:"signature"`
	Amount    float64 `json:"amount"`
	Nonce     uint64  `json:"nonce"`
	Timestamp int64   `json:"timestamp"`
	Type      TxType  `json:"type"`
}

// String returns the transaction as a string.
func (t Transaction) String() string {
	return fmt.Sprintf("%#v", t)
}

// Hash returns the hash of the transaction.
func (t Transaction) Hash() []byte {
	h := sha256.New()
	h.Write([]byte(t.String()))

	return h.Sum(nil)
}

// Verify verifies if the signature is valid.
func (t Transaction) Verify() error {
	// decode public key
	key, err := crypto.DecodePublicKey(util.HexDecode(t.Sender))
	if err != nil {
		return err
	}

	if !crypto.Verify(key, []byte("test"), util.HexDecode(t.Signature)) {
		return fmt.Errorf("%w: invalid signature", ErrInvalidTransaction)
	}

	//// check whether the signature is valid
	//if !crypto.Verify(key, []byte(fmt.Sprintf("%s%s%f", t.Sender, t.Receiver, t.Amount)), util.HexDecode(t.Signature)) {
	//	return fmt.Errorf("%w: invalid signature", ErrInvalidTransaction)
	//}

	return nil
}

// hashTransactions returns the hash of all given transactions.
func hashTransactions(transactions []Transaction) [][]byte {
	data := make([][]byte, 0, len(transactions))

	for _, t := range transactions {
		data = append(data, t.Hash())
	}

	return data
}
