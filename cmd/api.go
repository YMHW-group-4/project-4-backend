package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"backend/networking"
	"backend/util"
	"backend/wallet"

	"github.com/rs/cors"
	"github.com/rs/zerolog/log"
)

// API needs Node; main cannot be shared, thus refactoring needs to be done.
// Moving the API here until rewrite.
// Spoiler: refactoring will not be done.

var errInvalidHost = errors.New("invalid host")

// API represents the HTTP API.
type API struct {
	server *http.Server
	seed   string
	wg     sync.WaitGroup
}

// NewAPI creates a new HTTP API.
func NewAPI(port int, seed string) *API {
	mux := http.NewServeMux()

	mux.HandleFunc("/transaction", transaction)
	mux.HandleFunc("/freemoney", freeMoney)
	mux.HandleFunc("/wallets", wallets)
	mux.HandleFunc("/balance", balance)

	return &API{
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: cors.Default().Handler(mux),
		},
		seed: seed,
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

	a.Register()

	go func() {
		defer a.wg.Done()

		if err := a.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("api: failed to serve")
		}
	}()

	log.Info().Msg("api: running")
}

// getOutboundIP Get preferred outbound ip of this machine.
func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Debug().Err(err).Msg("api: failed to get outbound IP")
	}
	defer conn.Close()

	localAddr, ok := conn.LocalAddr().(*net.UDPAddr)
	if !ok {
		return net.IP{}
	}

	return localAddr.IP
}

// Register registers the node to the DNS seed.
func (a *API) Register() {
	if len(strings.TrimSpace(a.seed)) == 0 {
		log.Debug().Err(errInvalidHost).Msg("api: failed to register to DNS seed")

		return
	}

	url := fmt.Sprintf("%s/register_node?host=%s&port=%s", a.seed, getOutboundIP(), a.server.Addr)

	if _, err := http.Get(url); err != nil {
		log.Debug().Err(err).Msg("api: failed to register to DNS seed")

		return
	}

	log.Debug().Msg("api: registered to DNS seed")
}

// balance returns the balance of a wallet to a caller.
func balance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)

		return
	}

	sender := strings.TrimSpace(r.URL.Query().Get("sender"))

	if len(sender) == 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

		return
	}

	var b string

	account, err := node.blockchain.GetAccount(sender)
	if err != nil {
		b = "0.00"
	} else {
		b = account.Balance.String()
	}

	if err = json.NewEncoder(w).Encode(b); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	log.Debug().Str("endpoint", "balance").Msg("api: handled request")
}

// transaction creates and returns a new transaction to the caller.
func transaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)

		return
	}

	sender := strings.TrimSpace(r.URL.Query().Get("sender"))
	receiver := strings.TrimSpace(r.URL.Query().Get("receiver"))
	signature := strings.TrimSpace(r.URL.Query().Get("signature"))
	amount := strings.TrimSpace(r.URL.Query().Get("amount"))

	if len(sender) == 0 || len(receiver) == 0 || len(signature) == 0 || len(amount) == 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

		return
	}

	f, err := strconv.ParseFloat(amount, 64) //nolint
	if err != nil {
		http.Error(w, "parameter 'amount' invalid", http.StatusBadRequest)

		return
	}

	t, err := node.blockchain.CreateTransaction(sender, receiver, []byte(signature), f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	node.network.Publish(networking.Transaction, util.MarshalType(t))

	if err = json.NewEncoder(w).Encode(t); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	log.Debug().Str("endpoint", "transaction").Msg("api: handled request")
}

// freeMoney creates and returns a new transaction from genesis to the caller.
func freeMoney(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)

		return
	}

	sender := strings.TrimSpace(r.URL.Query().Get("sender"))
	amount := strings.TrimSpace(r.URL.Query().Get("amount"))

	if len(sender) == 0 || len(amount) == 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

		return
	}

	f, err := strconv.ParseFloat(amount, 64) //nolint
	if err != nil {
		http.Error(w, "parameter 'amount' invalid", http.StatusBadRequest)

		return
	}

	t, err := node.blockchain.CreateTransaction("genesis", sender, []byte(""), f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	node.network.Publish(networking.Transaction, util.MarshalType(t))

	if err = json.NewEncoder(w).Encode(t); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	log.Debug().Str("endpoint", "freemoney").Msg("api: handled request")
}

// wallets creates and returns a new wallet to the caller.
func wallets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)

		return
	}

	if err := json.NewEncoder(w).Encode(wallet.CreateWallet()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	log.Debug().Str("endpoint", "wallets").Msg("api: handled request")
}
