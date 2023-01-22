package blockchain

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MempoolTestSuite struct {
	suite.Suite
	mp   *mempool
	data []Transaction
}

func (suite *MempoolTestSuite) SetupTest() {
	suite.mp = newMempool()
}

func TestMempoolSuite(t *testing.T) {
	suite.Run(t, new(MempoolTestSuite))
}

func (suite *MempoolTestSuite) TestMempoolDoubleTransaction() {
	transaction := Transaction{
		Sender:    "Sender",
		Receiver:  "Receiver",
		Signature: string("Signature"),
		Amount:    10,
		Nonce:     1,
		Timestamp: time.Now().Unix(),
	}

	x, _ := json.MarshalIndent([]Transaction{transaction, transaction}, "", " ") //nolint
	log.Debug().Msgf("%s", x)

	err := suite.mp.add([]Transaction{transaction, transaction}...)

	assert.NotNil(suite.T(), err)
}

func TestMempoolRetrieveTransactions(t *testing.T) {
	mp := newMempool()

	_ = mp.add(transactions...)

	assert.Equal(t, 3, len(mp.retrieve(3)))
}
