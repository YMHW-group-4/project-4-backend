package wallet

import (
	"crypto/ecdsa"
	"crypto/rand"

	"github.com/ethereum/go-ethereum/crypto"
)

// Sign signs a hash using the private key.
func Sign(priv *ecdsa.PrivateKey, hash []byte) ([]byte, error) {
	return ecdsa.SignASN1(rand.Reader, priv, hash)
}

// Verify checks whether the signature is valid.
func Verify(pub *ecdsa.PublicKey, hash []byte, sig []byte) bool {
	return ecdsa.VerifyASN1(pub, hash, sig)
}

// EncodePublicKey encodes a public key.
func EncodePublicKey(key *ecdsa.PublicKey) []byte {
	return crypto.FromECDSAPub(key)
}

// EncodePrivateKey encodes a private key.
func EncodePrivateKey(key *ecdsa.PrivateKey) []byte {
	return crypto.FromECDSA(key)
}

// DecodePublicKey decodes a public key.
func DecodePublicKey(b []byte) (*ecdsa.PublicKey, error) {
	pub, err := crypto.UnmarshalPubkey(b)
	if err != nil {
		return nil, err
	}

	return pub, nil
}

// DecodePrivateKey decodes a private key.
func DecodePrivateKey(b []byte) (*ecdsa.PrivateKey, error) {
	priv, err := crypto.ToECDSA(b)
	if err != nil {
		return nil, err
	}

	return priv, nil
}
