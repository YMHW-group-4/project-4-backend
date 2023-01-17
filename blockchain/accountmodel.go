package blockchain

import (
	"sync"

	"backend/errors"
)

// account represents an account within the accountModel.
// It holds the balance and the number of transactions done by the account.
type account struct {
	balance      float32
	transactions uint32
}

// accountModel holds the accounts of all keys.
type accountModel struct {
	sync.RWMutex
	accounts map[string]*account
}

// newAccountModel creates a new accountModel.
func newAccountModel() *accountModel {
	return &accountModel{
		accounts: make(map[string]*account),
	}
}

// fromBlocks generates the balances from every block.
// Method should only be called on blockchain initialization.
// Assumption is that every block within the blockchain are valid.
// That also means that all transactions in a block are valid.
func (am *accountModel) fromBlocks(block ...Block) {
	var wg sync.WaitGroup

	for _, b := range block {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for _, transaction := range b.Transactions {
				am.Lock()

				rx := am.accounts[transaction.Sender]
				tx := am.accounts[transaction.Receiver]

				if rx != nil {
					rx.balance -= transaction.Amount
					rx.transactions++
				} else {
					// this should not happen.
					// sender should always exist; default balance is zero.
					// transaction should be verified before its being forged into a block.
					am.accounts[transaction.Sender] = &account{
						balance:      0,
						transactions: 1,
					}
				}

				if tx != nil {
					tx.balance += transaction.Amount
				} else {
					am.accounts[transaction.Receiver] = &account{
						balance:      transaction.Amount,
						transactions: 0,
					}
				}

				am.Unlock()
			}
		}()

		wg.Wait()
	}
}

// clear clears the accountModel.
func (am *accountModel) clear() {
	if len(am.accounts) > 0 {
		am.accounts = make(map[string]*account)
	}
}

// get returns the account associated with given key.
func (am *accountModel) get(key string) (*account, error) {
	if !am.exists(key) {
		return nil, errors.ErrInvalidInput("key does not exist")
	}

	am.RLock()
	defer am.RUnlock()

	return am.accounts[key], nil
}

// exists checks if the given key is already in the accountModel.
func (am *accountModel) exists(key string) bool {
	am.RLock()
	defer am.RUnlock()

	_, ok := am.accounts[key]

	return ok
}

// add adds the given key to the accountModel.
func (am *accountModel) add(key string, balance float32, transactions uint32) error {
	if am.exists(key) {
		return errors.ErrInvalidInput("key already exists")
	}

	// this should not happen
	if 0 > balance {
		return errors.ErrInvalidOperation("balance cannot be negative")
	}

	am.Lock()
	defer am.Unlock()

	am.accounts[key] = &account{
		balance:      balance,
		transactions: transactions,
	}

	return nil
}

// update updates the balance of the given key.
func (am *accountModel) update(key string, amount float32) error {
	if !am.exists(key) {
		return errors.ErrInvalidInput("key does not exist")
	}

	// this should not happen
	if 0 > (am.accounts[key].balance + amount) {
		return errors.ErrInvalidOperation("balance cannot be negative")
	}

	am.Lock()
	defer am.Unlock()

	am.accounts[key].balance += amount
	am.accounts[key].transactions++

	return nil
}
