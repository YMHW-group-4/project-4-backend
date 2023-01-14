package blockchain

import (
	"backend/errors"
	"github.com/rs/zerolog/log"
	"sync"
)

// AccountModel is the structure of the account model that is used in the blockchain.
type AccountModel struct {
	Balances map[string]float32
	wg       sync.WaitGroup
	rw       sync.RWMutex // TODO implement lock for concurrent map access
}

// NewAccountModel creates a new AccountModel.
func NewAccountModel() *AccountModel {
	return &AccountModel{
		Balances: make(map[string]float32),
	}
}

func (am *AccountModel) BalanceFromBlocks(blocks []Block) {
	if len(am.Balances) > 0 {
		am.Balances = make(map[string]float32)
	}

	for _, block := range blocks {
		am.wg.Add(1)

		go func(block Block) {
			defer am.wg.Done()
			// TODO
		}(block)
	}

	am.wg.Wait()
}

func (am *AccountModel) BalanceFromBlock(block Block) {
	for _, transaction := range block.Transactions {
		log.Debug().Msgf("%v", transaction)
		// TODO
	}
}

// Exists checks if the given key is already in the AccountModel.
func (am *AccountModel) Exists(key string) bool {
	_, ok := am.Balances[key]

	return ok
}

// Add adds the given key to the AccountModel.
func (am *AccountModel) Add(key string) error {
	if am.Exists(key) {
		return errors.ErrInvalidInput("key already exists")
	}

	am.Balances[key] = 0

	return nil
}

// Update updates the balance of the given key.
func (am *AccountModel) Update(key string, amount float32) error {
	if !am.Exists(key) {
		return errors.ErrInvalidInput("key does not exist")
	}

	// this should not happen
	if 0 > (am.Balances[key] + amount) {
		return errors.ErrInvalidOperation("balance cannot be negative")
	}

	am.Balances[key] += amount

	return nil
}
