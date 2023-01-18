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

	e := &account{
		balance:      0,
		transactions: 10000,
	}

	assert.Equal(t, e, am.accounts["genesis"])
}
