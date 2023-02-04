package wallet

import (
	"testing"

	"backend/util"

	"github.com/stretchr/testify/assert"
)

func TestEqualKeyPairs(t *testing.T) {
	m, _ := generateMnemonic(256)

	m1, priv1, pub1, _ := NewKeyPair(m, "")
	m2, priv2, pub2, _ := NewKeyPair(m, "")

	assert.Equal(t, m1, m2)
	assert.Equal(t, priv1, priv2)
	assert.Equal(t, pub1, pub2)
}

func TestSerializeKeys(t *testing.T) {
	m, _ := generateMnemonic(256)

	_, priv1, pub1, _ := NewKeyPair(m, "")

	privHexEncoded := util.HexEncode(EncodePrivateKey(priv1))
	pubHexEncoded := util.HexEncode(EncodePublicKey(pub1))

	privHexDecoded := util.HexDecode(privHexEncoded)
	pubHexDecoded := util.HexDecode(pubHexEncoded)

	priv2, _ := DecodePrivateKey(privHexDecoded)
	pub2, _ := DecodePublicKey(pubHexDecoded)

	assert.Equal(t, priv1, priv2)
	assert.Equal(t, pub1, pub2)
}
