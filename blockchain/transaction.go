package blockchain

type Transaction struct {
	PubKeyTx string
	PubKeyRx string
	Amount   float32
}

func CreateTransaction() Transaction {
	return Transaction{}
}
