package blockchain

type Block struct {
	Transactions []Transaction
	Hash         string
	Timestamp    int
}
