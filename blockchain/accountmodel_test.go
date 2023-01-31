package blockchain

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccountModelFromBlock(t *testing.T) {
	am := newAccountModel()

	_ = am.add("genesis", 10000, 0) //nolint

	var blocks []Block

	for i := 0; i < 100; i++ {
		block := Block{}

		for j := 0; j < 100; j++ {
			transaction := Transaction{
				Sender:   "genesis",
				Receiver: string(rune(rand.Int())),
				Amount:   1,
			}

			block.Transactions = append(block.Transactions, transaction)
		}

		blocks = append(blocks, block)
	}

	am.fromBlocks(blocks...)

	e := &Account{
		Balance:      ToCoin(0),
		Transactions: 10000,
	}

	assert.Equal(t, e.Transactions, am.accounts["genesis"].Transactions)
	assert.True(t, e.Balance.Equal(am.accounts["genesis"].Balance))
}

func TestAccountModelTransactions(t *testing.T) {
	am := newAccountModel()

	_ = am.add("genesis", 10000, 0)

	t1 := Transaction{Sender: "genesis", Receiver: "receiver", Amount: 20.15}
	t2 := Transaction{Sender: "genesis", Receiver: "receiver", Amount: 10.15}

	b := Block{Transactions: []Transaction{t1, t2}}

	am.fromBlocks(b)

	assert.True(t, ToCoin(30.30).Equal(am.accounts["receiver"].Balance))
}
