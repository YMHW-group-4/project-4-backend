package blockchain

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMempoolDoubleTransaction(t *testing.T) {
	mp := newMempool()

	transaction := Transaction{
		Sender:    "Sender",
		Receiver:  "Receiver",
		Signature: "Signature",
		Amount:    10,
		Nonce:     1,
		Timestamp: time.Now().Unix(),
	}

	err := mp.add([]Transaction{transaction, transaction}...)

	assert.NotNil(t, err)
}

func TestMempoolRetrieveTransactions(t *testing.T) {
	mp := newMempool()

	_ = mp.add(transactions...)

	assert.Equal(t, 3, len(mp.retrieve(3)))
}
