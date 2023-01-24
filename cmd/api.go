package main

import (
	"backend/networking"
	"backend/util"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/ripemd160"
	"net/http"
	"strconv"
	"sync"

	"github.com/rs/zerolog/log"
)

var SEED_HOST = "http://167.86.93.188"
var SEED_PORT = "3000"
var PORT = "8081"
var node *Node

// API needs Node; main cannot be shared, thus refactoring needs to be done.
// Moving the API here until rewrite.
// Spoiler: refactoring will not be done.

// API represents the HTTP API.
type API struct {
	server *http.Server
	wg     sync.WaitGroup
}

// NewAPI creates a new HTTP API.
func NewAPI(port int, n *Node) *API {
	mux := http.NewServeMux()

	mux.HandleFunc("/transaction", Transaction) //FIXME error checking
	mux.HandleFunc("/freemoney", Freemoney)     //FIXME error checking
	mux.HandleFunc("/wallets", Wallets)         // FIXME error checking
	mux.HandleFunc("/balance", Balance)         //TODO

	node = n

	return &API{
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		},
	}
}

// Stop stops the API.
func (a *API) Stop() {
	log.Info().Msg("api: shutting down")

	ctx := context.Background()

	defer ctx.Done()

	if err := a.server.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("api: failed to shutdown")
	}

	a.wg.Wait()

	log.Info().Msg("api: terminated")
}

// Start starts the API.
func (a *API) Start() {
	log.Info().Msg("api: starting")

	a.wg.Add(1)

	a.RegisterServer()

	go func() {
		defer a.wg.Done()

		if err := a.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("api: failed to serve")
		}
	}()

	log.Info().Msg("api: running")
}

func (a *API) RegisterServer() {
	host := "localhost"

	url := SEED_HOST + ":" + SEED_PORT + "/register_node?host=" + host + "&port=" + a.server.Addr
	fmt.Println(url)
	respose, err := http.Get(url)
	log.Debug().Msgf("%v", respose)
	if err != nil {
		log.Warn().Err(err).Msg("Error while registering to the seeder")
	}
}

func Transaction(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if req.Method == "POST" {

		sender := req.URL.Query().Get("sender")
		receiver := req.URL.Query().Get("receiver")
		signature := req.URL.Query().Get("signature")
		amount, err := strconv.ParseFloat(req.URL.Query().Get("amount"), 32)
		if err != nil {
			log.Error().Err(err).Msg("api: error")
		}

		_, err = node.blockchain.CreateTransaction(sender, receiver, []byte(signature), float32(amount))
		if err != nil {
			log.Error().Err(err).Msg("error:")
		} else {
			log.Debug().Msg("nice one")
		}

		json.NewEncoder(w).Encode("{data: \"Test\", lala: \"ytes\"}")

	}
}

func Freemoney(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	sender := req.URL.Query().Get("sender")

	log.Debug().Str("sender", sender).Msg("api: ")

	amount, err := strconv.ParseFloat(req.URL.Query().Get("amount"), 32)
	if err != nil {
		log.Error().Err(err).Msg("api: error")
	}

	t, err := node.blockchain.CreateTransaction("genesis", sender, []byte(""), float32(amount))
	if err != nil {
		log.Error().Err(err).Msg("api: error")
	}

	node.network.Publish(networking.Transaction, util.MarshalType(t))

	log.Debug().Msgf("%v", t)

}

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
	hasher := ripemd160.New()
	_, _ = hasher.Write(hashedPublicKey[:])
	publicRipeMd := hasher.Sum(nil)
	return publicRipeMd
}

func Wallets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		data := CreateWallet()
		json.NewEncoder(w).Encode(data)
	}
}

func Balance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		json.NewEncoder(w).Encode("Nee")
	}
}
