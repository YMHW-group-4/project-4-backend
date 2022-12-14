package blockchain

import "time"

type Transaction struct {
	PubKeySender	string
	PubKeyReceiver	string
	Amount 			float32
	Id				string
	Timestamp		time.Time
}

func createTransaction(sender string, receiver string, amount float32, id string) Transaction{
	return Transaction {
		PubKeySender: sender,
		PubKeyReceiver: receiver,
		Amount: amount,
		Id: id,
		Timestamp: time.Now(),
	}
}
