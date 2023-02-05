package crypto

import (
	"crypto/ecdsa"
	"embed"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
)

//go:embed keys
var keys embed.FS

// generateGenesis generates a new keypair for genesis.
func generateGenesis() error {
	priv, err := crypto.GenerateKey()
	if err != nil {
		return err
	}

	if err = os.WriteFile("keys/genesis.priv.pem", PemEncodePrivateKey(priv), 0o600); err != nil {
		return err
	}

	if err = os.WriteFile("keys/genesis.pub.pem", PemEncodePublicKey(&priv.PublicKey), 0o600); err != nil {
		return err
	}

	return nil
}

// Genesis returns the private and public key of genesis.
func Genesis() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	priv, err := keys.ReadFile("keys/genesis.priv.pem")
	if err != nil {
		return nil, nil, err
	}

	pub, err := keys.ReadFile("keys/genesis.pub.pem")
	if err != nil {
		return nil, nil, err
	}

	privKey, err := PemDecodePrivateKey(priv)
	if err != nil {
		return nil, nil, err
	}

	pubKey, err := PemDecodePublicKey(pub)
	if err != nil {
		return nil, nil, err
	}

	return privKey, pubKey, nil
}
