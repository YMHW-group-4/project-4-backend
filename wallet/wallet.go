package wallet

import (
	"backend/crypto"
	"backend/util"
)

// A wallet should not be generated on the Node.
// Wallet should be seen as a standalone application.
// Due to time constraints, a wallet is generated on the node via its API.

// Wallet represents the private and public key within the blockchain.
type Wallet struct {
	Mnemonic string
	Priv     string
	Pub      string
}

// CreateWallet creates a new Wallet.
func CreateWallet(mnemonic string, password string) (*Wallet, error) {
	m, priv, pub, err := newKeyPair(mnemonic, password)
	if err != nil {
		return nil, err
	}

	return &Wallet{
		Mnemonic: m,
		Priv:     util.HexEncode(crypto.EncodePrivateKey(priv)),
		Pub:      util.HexEncode(crypto.EncodePublicKey(pub)),
	}, nil
}
