package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"backend/api/endpoints"

	"github.com/rs/zerolog/log"
)

// API represents the HTTP API.
type API struct {
	server *http.Server
	wg     sync.WaitGroup
}

// NewAPI creates a new HTTP API.
func NewAPI(port int) *API {
	mux := http.NewServeMux()

	mux.HandleFunc("/transaction", endpoints.Transaction)
	mux.HandleFunc("/wallet", endpoints.Wallet)

	return &API{
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		},
	}
}

// Stop stops the API.
func (a *API) Stop() error {
	ctx := context.Background()

	defer ctx.Done()

	if err := a.server.Shutdown(ctx); err != nil {
		return err
	}

	a.wg.Wait()

	return nil
}

// Start starts the API.
func (a *API) Start() {
	a.wg.Add(1)

	go func() {
		defer a.wg.Done()

		if err := a.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Error().Err(err).Msg("API: failed to serve")
		}
	}()
}
