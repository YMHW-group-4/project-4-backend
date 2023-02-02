package wallet

import "github.com/tyler-smith/go-bip39"

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
