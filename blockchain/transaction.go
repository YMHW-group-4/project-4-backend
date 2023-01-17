package blockchain

type Transaction struct {
	Sender    string
	Receiver  string
	Signature string
	Amount    float32
	Nonce     uint32
	Timestamp int64
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
