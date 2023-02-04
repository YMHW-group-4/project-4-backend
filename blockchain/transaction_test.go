package blockchain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TransactionTestSuite struct {
	suite.Suite
	priv *ecdsa.PrivateKey
	pub  *ecdsa.PublicKey
}

func (suite *TransactionTestSuite) SetupTest() {
	suite.priv, _ = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	suite.pub = &suite.priv.PublicKey
}

func TestTransactionTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionTestSuite))
}

func (suite *TransactionTestSuite) TestTransactionSignature() {
	sig, _ := ecdsa.SignASN1(rand.Reader, suite.priv, []byte("signature"))

	assert.True(suite.T(), ecdsa.VerifyASN1(suite.pub, []byte("signature"), sig))
}
