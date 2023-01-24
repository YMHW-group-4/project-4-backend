package main

import (
	"backend/networking"
	"backend/util"
	"backend/wallet"
	"context"

	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/rs/cors"
	"github.com/rs/zerolog/log"
)

var SeedHost = "http://167.86.93.188"
var SeedPort = 3000

// API needs Node; main cannot be shared, thus refactoring needs to be done.
// Moving the API here until rewrite.
// Spoiler: refactoring will not be done.

// API represents the HTTP API.
type API struct {
	server *http.Server
	wg     sync.WaitGroup
}

// NewAPI creates a new HTTP API.
func NewAPI(port int) *API {
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

func (a *API) Register() {
	host := "localhost"

	url := fmt.Sprintf("%s:%d/register_node?host=%s%s", SeedHost, SeedPort, host, a.server.Addr)

	if _, err := http.Get(url); err != nil {
		log.Debug().Err(err).Msg("api: failed to register")

		return
	}

	log.Debug().Msg("api: registered to seed")
}

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

	var b float32

	account, err := node.blockchain.GetAccount(sender)
	if err != nil {
		b = 0
	} else {
		b = account.Balance
	}

	if err = json.NewEncoder(w).Encode(b); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	log.Debug().Str("endpoint", "balance").Msg("api: handled request")
}

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

	f, err := strconv.ParseFloat(amount, 32)
	if err != nil {
		http.Error(w, "parameter 'amount' invalid", http.StatusBadRequest)

		return
	}

	t, err := node.blockchain.CreateTransaction(sender, receiver, []byte(signature), float32(f))
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

	f, err := strconv.ParseFloat(amount, 32)
	if err != nil {
		http.Error(w, "parameter 'amount' invalid", http.StatusBadRequest)

		return
	}

	t, err := node.blockchain.CreateTransaction("genesis", sender, []byte(""), float32(f))
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
