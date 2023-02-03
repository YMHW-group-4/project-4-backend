package wallet

import (
	"crypto/ecdsa"
	"encoding/pem"
	"errors"
)

// errPemInvalidBlock is the error when a PEM block is invalid.
var errPemInvalidBlock = errors.New("invalid block")

// pemDecodePublicKey decodes a PEM-encoded ECDSA public key.
func pemDecodePublicKey(key []byte) (*ecdsa.PublicKey, error) {
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

// pemDecodePrivateKey decodes a PEM-encoded ECDSA private key.
func pemDecodePrivateKey(key []byte) (*ecdsa.PrivateKey, error) {
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

// pemEncodePublicKey encodes an ECDSA public key to PEM format.
func pemEncodePublicKey(key *ecdsa.PublicKey) []byte {
	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: EncodePublicKey(key),
	}

	return pem.EncodeToMemory(block)
}

// pemEncodePrivateKey encodes an ECDSA private key to PEM format.
func pemEncodePrivateKey(key *ecdsa.PrivateKey) []byte {
	keyBlock := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: EncodePrivateKey(key),
	}

	return pem.EncodeToMemory(keyBlock)
}
