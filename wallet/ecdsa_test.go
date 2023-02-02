package wallet

import (
	"backend/util"
	"crypto/elliptic"
	"fmt"
	"github.com/rs/zerolog/log"
	"testing"
)

func TestKey(t *testing.T) {
	priv, pub, _ := generateKey(elliptic.P256())

	privEncoded, err := EncodePrivateKey(priv)
	if err != nil {
		log.Error().Err(err).Msg("EncodePrivateKey")
	}

	privHex := util.HexEncode(privEncoded)

	pubEncoded, err := EncodePublicKey(pub)
	if err != nil {
		log.Error().Err(err).Msg("EncodePublicKey")
	}

	pubHex := util.HexEncode(pubEncoded)

	privDecoded, err := DecodePrivateKey(util.HexDecode(privHex))
	if err != nil {
		log.Error().Err(err).Msg("DecodePrivateKey")
	}

	pubDecoded, err := DecodePublicKey(util.HexDecode(pubHex))
	if err != nil {
		log.Error().Err(err).Msg("DecodePublicKey")
	}

	hash := []byte(fmt.Sprintf("allo"))

	s1, _ := Sign(priv, hash)
	s2, _ := Sign(privDecoded, hash)

	if Verify(pubDecoded, hash, s1) {
		log.Debug().Msg("s1 pubDecoded: true")
	}

	if Verify(pub, hash, s1) {
		log.Debug().Msg("s1 pub: true")
	}

	if Verify(pubDecoded, hash, s2) {
		log.Debug().Msg("s2 pubDecoded: true")
	}

	if Verify(pubDecoded, hash, s2) {
		log.Debug().Msg("s2 pub: true")
	}
}
