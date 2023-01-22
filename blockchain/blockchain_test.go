package blockchain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"github.com/rs/zerolog/log"
	"testing"

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
	priv, _ := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	pub := priv.PublicKey

	_ = suite.b.am.add(fmt.Sprintf("%v", pub), 20, 0)
	hash := []byte(fmt.Sprintf("%s%s%d", pub, "receiver", 20))

	//
	sig, _ := ecdsa.SignASN1(rand.Reader, priv, hash)
	//
	_, err := suite.b.CreateTransaction(fmt.Sprintf("%v", pub), "receiver", sig, 20)
	log.Debug().Err(err).Send()

	log.Debug().Msgf("%v", pub)

}
