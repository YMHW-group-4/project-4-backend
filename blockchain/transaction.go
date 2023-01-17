package blockchain

// Transaction represents a transaction within the blockchain.
type Transaction struct {
	Sender    string  `json:"sender"`
	Receiver  string  `json:"receiver"`
	Signature string  `json:"signature"`
	Amount    float32 `json:"amount"`
	Nonce     uint32  `json:"nonce"`
	Timestamp int64   `json:"timestamp"`
}

//func CreateTransaction(amount float32) Transaction {
//	return Transaction{
//
//		Amount:    amount,
//		Nonce:
//		Timestamp: time.Now().Unix(),
//	}
//}

func (t Transaction) validate() {

}
