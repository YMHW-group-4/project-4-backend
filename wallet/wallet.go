package wallet

import (
	"crypto/ecdsa"
	"embed"
	"os"

	"backend/util"
)

//go:embed keys
var keys embed.FS

// Wallet represents the private and public key within the blockchain.
type Wallet struct {
	Mnemonic string
	Priv     *ecdsa.PrivateKey
	Pub      *ecdsa.PublicKey
}

// CreateWallet creates a new Wallet.
func CreateWallet(mnemonic string, password string) (*Wallet, error) {
	m, priv, pub, err := NewKeyPair(mnemonic, password)
	if err != nil {
		return nil, err
	}

	return &Wallet{
		Mnemonic: m,
		Priv:     priv,
		Pub:      pub,
	}, nil
}

// Public returns the public key formatted as a string.
func (w *Wallet) Public() string {
	return util.HexEncode(EncodePublicKey(w.Pub))
}

// Private returns the private key formatted as a string.
func (w *Wallet) Private() string {
	return util.HexEncode(EncodePrivateKey(w.Priv))
}

// generateGenesis generates a new keypair for genesis.
func generateGenesis() error {
	w, err := CreateWallet("", "genesis")
	if err != nil {
		return err
	}

	if err = os.WriteFile("keys/genesis.priv.pem", pemEncodePrivateKey(w.Priv), 0600); err != nil { //nolint
		return err
	}

	if err = os.WriteFile("keys/genesis.pub.pem", pemEncodePublicKey(w.Pub), 0600); err != nil { //nolint
		return err
	}

	return nil
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
