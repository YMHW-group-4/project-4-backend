package blockchain

import (
	"fmt"
	"testing"

	"backend/crypto"
	"backend/util"
	"backend/wallet"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BlockchainTestSuite struct {
	suite.Suite
	b *Blockchain
}

func (suite *BlockchainTestSuite) SetupTest() {
	suite.b = NewBlockchain()
}

func TestBlockchainTestSuite(t *testing.T) {
	suite.Run(t, new(BlockchainTestSuite))
}

func (suite *BlockchainTestSuite) TestBlockchainCreateTransaction() {
	_, priv, pub, _ := wallet.NewKeyPair("", "")

	sender := util.HexEncode(crypto.EncodePublicKey(pub))
	receiver := "receiver"

	_ = suite.b.am.add(sender, 20)

	hash := []byte(fmt.Sprintf("%s%s%f", sender, receiver, float64(20)))

	sig, _ := crypto.Sign(priv, hash)

	_, err := suite.b.CreateTransaction(sender, receiver, sig, 20)

	assert.Nil(suite.T(), err)
}
