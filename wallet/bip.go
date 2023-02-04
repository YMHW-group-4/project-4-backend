package wallet

import (
	"crypto/ecdsa"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mr-tron/base58"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

// generateMnemonic generates a new mnenomic according to the bip39 specification.
func generateMnemonic(bitSize int) (string, error) {
	entropy, err := bip39.NewEntropy(bitSize)
	if err != nil {
		return "", err
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}

	return mnemonic, nil
}

// generateSeed generates a new seed from a mnemonic and a password.
func generateSeed(mnemonic string, password string) ([]byte, error) {
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, bip39.ErrInvalidMnemonic
	}

	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, password)
	if err != nil {
		return nil, err
	}

	return seed, nil
}

// newMasterKey creates a new master key according to the bip32 specification.
func newMasterKey(seed []byte) (*bip32.Key, error) {
	return bip32.NewMasterKey(seed)
}

// deriveECDSA derives the ECDSA keys according to the bip44 specification from a bip32 master key.
// source: https://gist.github.com/miguelmota/f56fa0b01e8c6c649a6c4f0ee7337aab
func deriveECDSA(master *bip32.Key) (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	// derive ECDSA according to the bip44 specification (m/44)
	key, err := master.NewChildKey(2147483648 + 44)
	if err != nil {
		return nil, nil, err
	}

	decoded, err := base58.Decode(key.B58Serialize())
	if err != nil {
		return nil, nil, err
	}

	priv := decoded[46:78]

	privECDSA, err := crypto.ToECDSA(priv)
	if err != nil {
		return nil, nil, err
	}

	return privECDSA, &privECDSA.PublicKey, nil
}

// NewKeyPair creates a new keypair.
func NewKeyPair(mnemonic string, password string) (string, *ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	m := mnemonic

	var err error

	if len(strings.TrimSpace(m)) == 0 {
		m, err = generateMnemonic(256)
		if err != nil {
			return "", nil, nil, err
		}
	}

	seed, err := generateSeed(m, password)
	if err != nil {
		return "", nil, nil, err
	}

	key, err := newMasterKey(seed)
	if err != nil {
		return "", nil, nil, err
	}

	priv, pub, err := deriveECDSA(key)
	if err != nil {
		return "", nil, nil, err
	}

	return mnemonic, priv, pub, nil
}
