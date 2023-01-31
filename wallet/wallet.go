package wallet

import (
	"backend/util"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
)

// Wallet represents a wallet in the context of a blockchain.
const (
	checksumLength = 4
	walletVersion  = byte(0x00)
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func (w *Wallet) Address() []byte {
	pubHash := PublicKeyHash(w.PublicKey)
	versionedHash := append([]byte{walletVersion}, pubHash...)
	checksum := Checksum(versionedHash)
	finalHash := append(versionedHash, checksum...)

	return util.Base58Encode(finalHash)
}

func CreateWallet() map[string]any {
	wallet := MakeWallet()
	address := wallet.Address()

	body := make(map[string]any)
	body["private"] = address
	body["public"] = wallet.PublicKey

	return body
}

func MakeWallet() *Wallet {
	privateKey, publicKey := NewKeyPair()

	return &Wallet{privateKey, publicKey}
}

func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, _ := ecdsa.GenerateKey(curve, rand.Reader)
	pub := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pub
}

func Checksum(ripeMdHash []byte) []byte {
	firstHash := sha256.Sum256(ripeMdHash)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:checksumLength]
}

func PublicKeyHash(publicKey []byte) []byte {
	hashedPublicKey := sha256.Sum256(publicKey)
	hasher := sha256.New()
	_, _ = hasher.Write(hashedPublicKey[:])
	publicRipeMd := hasher.Sum(nil)

	return publicRipeMd
}
