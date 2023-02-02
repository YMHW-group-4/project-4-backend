package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"errors"
)

// errTypeAssertionFailed is the error when a type assertion cannot be made.
var errTypeAssertionFailed = errors.New("type assertion failed")

// generateKey generates a new keypair.
func generateKey(curve elliptic.Curve) (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	priv, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	return priv, &priv.PublicKey, nil
}

// Sign signs a hash using the private key.
func Sign(priv *ecdsa.PrivateKey, hash []byte) ([]byte, error) {
	return ecdsa.SignASN1(rand.Reader, priv, hash)
}

// Verify checks whether the signature is valid.
func Verify(pub *ecdsa.PublicKey, hash []byte, sig []byte) bool {
	return ecdsa.VerifyASN1(pub, hash, sig)
}

// EncodePublicKey encodes a public key.
func EncodePublicKey(key *ecdsa.PublicKey) ([]byte, error) {
	return x509.MarshalPKIXPublicKey(key)
}

// EncodePrivateKey encodes a private key.
func EncodePrivateKey(key *ecdsa.PrivateKey) ([]byte, error) {
	return x509.MarshalECPrivateKey(key)
}

// DecodePublicKey decodes a public key.
func DecodePublicKey(b []byte) (*ecdsa.PublicKey, error) {
	pub, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		return nil, err
	}

	ecdsaPub, ok := pub.(*ecdsa.PublicKey)
	if !ok {
		return nil, errTypeAssertionFailed
	}

	return ecdsaPub, nil
}

// DecodePrivateKey decodes a private key.
func DecodePrivateKey(b []byte) (*ecdsa.PrivateKey, error) {
	priv, err := x509.ParseECPrivateKey(b)
	if err != nil {
		return nil, err
	}

	return priv, nil
}
