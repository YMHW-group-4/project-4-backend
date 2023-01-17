package blockchain

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
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

	assert.Equal(t, uint32(10000), am.accounts["genesis"].transactions)
	assert.Equal(t, float32(0), am.accounts["genesis"].balance)
}
