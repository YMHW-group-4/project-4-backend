package wallet

import (
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

func NewMasterKey(mnemonic string, password string) (*bip32.Key, error) {
	seed := bip39.NewSeed(mnemonic, password)

	key, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, err
	}

	return key, nil

}

func NewChildKey(key bip32.Key) {
	key.NewChildKey()
}
