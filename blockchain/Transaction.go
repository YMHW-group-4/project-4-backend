package blockchain

import (
	"crypto"
	"encoding/json"
)

type Transaction struct {
	amount int16
	payer  crypto.PublicKey
	payee  crypto.PublicKey
}

func CreateTransaction(amount int16, payer crypto.PublicKey, payee crypto.PublicKey) Transaction {
	return Transaction{
		amount: amount,
		payer:  payer,
		payee:  payee,
	}
}

func (transaction *Transaction) ToString() string {
	data, _ := json.Marshal(transaction)
	return string(data)
}
