package wallet

import (
	"backend/errors"

	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

// generateMasterKey generates a new private master key from a mnemonic and a password.
func generateMasterKey(mnemonic string, password string) (*bip32.Key, error) {
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, bip39.ErrInvalidMnemonic
	}

	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, password)
	if err != nil {
		return nil, err
	}

	key, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, err
	}

	return key, nil
}

// generateChildKey generates a new child key from given master key.
func generateChildKey(key *bip32.Key, childID uint32) (*bip32.Key, error) {
	if !key.IsPrivate {
		return nil, errors.ErrInvalidInput("non-private key provided")
	}

	child, err := key.NewChildKey(childID)
	if err != nil {
		return nil, err
	}

	return child, nil
}
