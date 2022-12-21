package blockchain

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
)

type Wallet struct {
	publicKey  crypto.PublicKey
	PrivateKey *rsa.PrivateKey
	balance    int64
}

func CreateWallet() Wallet {
	priv, pub := getKeyPain()
	return Wallet{
		publicKey:  pub,
		PrivateKey: priv,
		balance:    0,
	}

}

func (wallet Wallet) SendMoney(amount int16, payeePublicKey crypto.PublicKey, chain *Chain) {
	newTransaction := Transaction{
		amount: amount,
		payer:  wallet.publicKey,
		payee:  payeePublicKey,
	}
	// TODO: Create a signature with the private and public keys.
	// TODO: Chang the private key to a signature
	// TODO: change the return to the singleton Chain
	chain.AddBlock(newTransaction, wallet.publicKey, wallet.PrivateKey)
}

func getKeyPain() (*rsa.PrivateKey, crypto.PublicKey) {
	bitSize := 64
	priv, _ := rsa.GenerateKey(rand.Reader, bitSize)
	pub := priv.Public()
	return priv, pub
}
