package blockchain

import (
	"fmt"
	"sync"

	"backend/errors"
)

// mempool represents the memory pool within the Blockchain.
// It acts as a storage for unconfirmed Transactions.
type mempool struct {
	sync.RWMutex
	pool map[string]Transaction
}

// newMempool creates a new memory pool.
func newMempool() *mempool {
	return &mempool{
		pool: make(map[string]Transaction),
	}
}

// clear clears the internal pool.
func (mp *mempool) clear() {
	mp.Lock()
	defer mp.Unlock()

	mp.pool = make(map[string]Transaction)
}

// exists checks if the given key is already in the mempool.
func (mp *mempool) exists(key string) bool {
	mp.RLock()
	defer mp.RUnlock()

	_, ok := mp.pool[key]

	return ok
}

// add adds a transaction to the mempool.
func (mp *mempool) add(transaction ...Transaction) error {
	var err error

	for _, tx := range transaction {
		key := tx.string()

		if mp.exists(key) {
			err = errors.ErrInvalidOperation(fmt.Sprintf("key %s already exists", key))

			continue
		}

		mp.Lock()

		mp.pool[key] = tx

		mp.Unlock()
	}

	return err
}

// retrieve retrieves transactions from the mempool.
// If an amount of zero is passed, all transactions in the mempool will be returned.
func (mp *mempool) retrieve(amount uint16) []Transaction {
	mp.RLock()
	defer mp.RUnlock()

	if amount == 0 {
		amount = uint16(len(mp.pool))
	}

	transactions := make([]Transaction, 0, amount)

	for _, transaction := range mp.pool {
		if uint16(len(transactions)) == amount {
			break
		}

		transactions = append(transactions, transaction)
	}

	return transactions
}

// delete removes a transaction from the mempool.
func (mp *mempool) delete(transaction ...Transaction) error {
	var err error

	for _, tx := range transaction {
		key := tx.string()

		if !mp.exists(key) {
			err = errors.ErrInvalidOperation(fmt.Sprintf("key %s does not exists", key))

			continue
		}

		mp.Lock()

		delete(mp.pool, key)

		mp.Unlock()
	}

	return err
}
