package blockchain

import (
	"sync"

	"backend/errors"
)

// accountModel holds the balances of all keys.
type accountModel struct {
	balances map[string]float32
	mu       sync.RWMutex
}

// newAccountModel creates a new accountModel.
func newAccountModel() *accountModel {
	return &accountModel{
		balances: make(map[string]float32),
	}
}

// fromBlocks generates the balances from every block.
// Should only be called on blockchain boot. Newly created blocks to be
// retrieved by fromBlock.
func (am *accountModel) fromBlocks(blocks []Block) {
	var wg sync.WaitGroup

	if len(am.balances) > 0 {
		am.balances = make(map[string]float32)
	}

	for _, block := range blocks {
		wg.Add(1)

		go func(block Block) {
			defer wg.Done()

			am.fromBlock(block)
		}(block)
	}

	wg.Wait()
}

// fromBlock generates the balances from a single block.
func (am *accountModel) fromBlock(block Block) {
	// FIXME could be changed; transaction structure might be changed.
	for _, transaction := range block.Transactions {
		am.mu.Lock()

		am.balances[transaction.PubKeyTx] = am.balances[transaction.PubKeyTx] - transaction.Amount
		am.balances[transaction.PubKeyRx] = am.balances[transaction.PubKeyRx] + transaction.Amount

		am.mu.Unlock()
	}
}

// exists checks if the given key is already in the accountModel.
func (am *accountModel) exists(key string) bool {
	am.mu.RLock()
	defer am.mu.RUnlock()

	_, ok := am.balances[key]

	return ok
}

// add adds the given key to the accountModel.
func (am *accountModel) add(key string) error {
	if am.exists(key) {
		return errors.ErrInvalidInput("key already exists")
	}

	am.mu.Lock()
	defer am.mu.Unlock()

	am.balances[key] = 0

	return nil
}

// update updates the balance of the given key.
func (am *accountModel) update(key string, amount float32) error {
	if !am.exists(key) {
		return errors.ErrInvalidInput("key does not exist")
	}

	// this should not happen
	if 0 > (am.balances[key] + amount) {
		return errors.ErrInvalidOperation("balance cannot be negative")
	}

	am.mu.Lock()
	defer am.mu.Unlock()

	am.balances[key] += amount

	return nil
}
