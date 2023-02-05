package wallet

import (
	"crypto/ecdsa"
	"embed"
	"os"

	"backend/crypto"
	"backend/util"
)

// A wallet should not be generated on the Node.
// Wallet should be seen as a standalone application.
// Due to time constraints, a wallet is generated on the node via its API.

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
	return util.HexEncode(crypto.EncodePublicKey(w.Pub))
}

// Private returns the private key formatted as a string.
func (w *Wallet) Private() string {
	return util.HexEncode(crypto.EncodePrivateKey(w.Priv))
}

// generateGenesis generates a new keypair for genesis.
func generateGenesis() error {
	w, err := CreateWallet("", "genesis")
	if err != nil {
		return err
	}

	if err = os.WriteFile("keys/genesis.priv.pem", crypto.PemEncodePrivateKey(w.Priv), 0o600); err != nil {
		return err
	}

	if err = os.WriteFile("keys/genesis.pub.pem", crypto.PemEncodePublicKey(w.Pub), 0o600); err != nil {
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

	privKey, err := crypto.PemDecodePrivateKey(priv)
	if err != nil {
		return nil, err
	}

	pubKey, err := crypto.PemDecodePublicKey(pub)
	if err != nil {
		return nil, err
	}

	return &Wallet{
		Priv: privKey,
		Pub:  pubKey,
	}, nil
}
