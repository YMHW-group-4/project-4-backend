package blockchain

// Mempool represents the memory pool within the Blockchain.
// It acts as a storage for unconfirmed Transactions.
type Mempool struct {
	Transactions []Transaction
}

func NewMempool() *Mempool {
	return &Mempool{
		Transactions: make([]Transaction, 0),
	}
}
