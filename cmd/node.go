package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend/blockchain"
	"backend/networking"

	"github.com/rs/zerolog/log"
)

// version is interpolated during build time.
var version string

// node represents a singular blockchain node.
type node struct {
	uptime     time.Time
	network    *networking.Network
	blockchain *blockchain.Blockchain
	api        any
}

// newNode creates a new node with given configuration.
func newNode(config configuration) (*node, error) {
	network, err := networking.NewNetwork(config.port)
	if err != nil {
		return nil, err
	}

	return &node{
		network:    network,
		blockchain: nil,
		api:        nil,
	}, nil
}

// run starts all services required by the node.
func (node *node) run() {
	if err := node.network.Start(); err != nil {
		log.Fatal().Err(err).Msg("node: failed to run")
	}

	node.uptime = time.Now()
}

// handleSigterm executes when termination from operating system is received.
func (node *node) handleSigterm() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)

	<-sigc

	log.Info().Msg("node: shutting down")

	if err := node.network.Close(); err != nil {
		log.Error().Err(err).Msg("node: failed to close network")
	}

	log.Info().Msg("node: terminated")
}
