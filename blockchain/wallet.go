package blockchain

import (
	"math/rand"
)

type Wallet struct {
	privateKey string
	publicKey  string
}

const keyBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func randomStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = keyBytes[rand.Intn(len(keyBytes))]
	}
	return string(b)
}

// func genPubKey() string {

// }

// func genPrivKey() string {

// }

// func createWallet(privKey string, publicKey string) Wallet {
// 	return Wallet()
// }
