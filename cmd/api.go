package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/rs/zerolog/log"
)

// API needs Node; main cannot be shared, thus refactoring needs to be done.
// Moving the API here until rewrite.
// Spoiler: refactoring will not be done.

// API represents the HTTP API.
type API struct {
	server *http.Server
	wg     sync.WaitGroup
	node   *Node
}

// NewAPI creates a new HTTP API.
func NewAPI(port int, node *Node) *API {
	mux := http.NewServeMux()

	// TODO refactor this
	mux.HandleFunc("/transaction", Transaction)

	return &API{
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		},
		node: node,
	}
}
tus
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

	go func() {
		defer a.wg.Done()

		if err := a.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("api: failed to serve")
		}
	}()

	log.Info().Msg("api: running")
}

func Transaction(w http.ResponseWriter, req *http.Request) {

}
