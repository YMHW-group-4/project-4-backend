package crypto

import (
	"crypto/ecdsa"
	"encoding/pem"
	"errors"
)

// errPemInvalidBlock is the error when a PEM block is invalid.
var errPemInvalidBlock = errors.New("invalid block")

// PemDecodePublicKey decodes a PEM-encoded ECDSA public key.
func PemDecodePublicKey(key []byte) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode(key)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errPemInvalidBlock
	}

	pub, err := DecodePublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return pub, nil
}

// PemDecodePrivateKey decodes a PEM-encoded ECDSA private key.
func PemDecodePrivateKey(key []byte) (*ecdsa.PrivateKey, error) {
	var block *pem.Block

	for {
		block, key = pem.Decode(key)

		if block == nil {
			return nil, errPemInvalidBlock
		}

		if block.Type == "EC PRIVATE KEY" {
			break
		}
	}

	priv, err := DecodePrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}

// PemEncodePublicKey encodes an ECDSA public key to PEM format.
func PemEncodePublicKey(key *ecdsa.PublicKey) []byte {
	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: EncodePublicKey(key),
	}

	return pem.EncodeToMemory(block)
}

// PemEncodePrivateKey encodes an ECDSA private key to PEM format.
func PemEncodePrivateKey(key *ecdsa.PrivateKey) []byte {
	keyBlock := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: EncodePrivateKey(key),
	}

	return pem.EncodeToMemory(keyBlock)
}
