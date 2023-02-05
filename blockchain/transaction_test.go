package blockchain

import (
	"crypto/ecdsa"
	"testing"

	"backend/crypto"
	"backend/wallet"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TransactionTestSuite struct {
	suite.Suite
	priv *ecdsa.PrivateKey
	pub  *ecdsa.PublicKey
}

func (suite *TransactionTestSuite) SetupTest() {
	_, priv, pub, _ := wallet.NewKeyPair("", "")
	suite.priv = priv
	suite.pub = pub
}

func TestTransactionTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionTestSuite))
}

func (suite *TransactionTestSuite) TestTransactionSignature() {
	sig, _ := crypto.Sign(suite.priv, []byte("signature"))

	assert.True(suite.T(), crypto.Verify(suite.pub, []byte("signature"), sig))
}
