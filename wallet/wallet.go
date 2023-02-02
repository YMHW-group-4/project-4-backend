package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"embed"

	"backend/util"
)

//go:embed keys
var keys embed.FS

// Wallet represents the private and public key within the blockchain.
type Wallet struct {
	Priv *ecdsa.PrivateKey
	Pub  *ecdsa.PublicKey
}

// CreateWallet creates a new Wallet.
func CreateWallet() (*Wallet, error) {
	priv, pub, err := generateKey(elliptic.P256())
	if err != nil {
		return nil, err
	}

	return &Wallet{
		Priv: priv,
		Pub:  pub,
	}, nil
}

// Public returns the public key formatted as a string.
func (w *Wallet) Public() string {
	pub, err := EncodePublicKey(w.Pub)
	if err != nil {
		return ""
	}

	return util.HexEncode(pub)
}

// Private returns the private key formatted as a string.
func (w *Wallet) Private() string {
	priv, err := EncodePrivateKey(w.Priv)
	if err != nil {
		return ""
	}

	return util.HexEncode(priv)
}

// GenesisWallet returns the Wallet of genesis.
func GenesisWallet() (*Wallet, error) {
	priv, err := keys.ReadFile("keys/genesis.priv.pem")
	if err != nil {
		return nil, err
	}

	pub, err := keys.ReadFile("keys/genesis.pub.pem")
	if err != nil {
		return nil, err
	}

	privKey, err := pemDecodePrivateKey(priv)
	if err != nil {
		return nil, err
	}

	pubKey, err := pemDecodePublicKey(pub)
	if err != nil {
		return nil, err
	}

	return &Wallet{
		Priv: privKey,
		Pub:  pubKey,
	}, nil
}
