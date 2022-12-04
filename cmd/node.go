package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
)

// version is interpolated during build time.
var version string

// node represents a singular blockchain node.
type node struct {
	uptime     time.Time
	network    any
	blockchain any
	api        any
}

// newNode creates a new node with given configuration.
func newNode(config configuration) (*node, error) {
	return &node{
		network:    nil,
		blockchain: nil,
		api:        nil,
	}, nil
}

// run starts all services required by the node.
func (node *node) run() {
	node.uptime = time.Now()
}

// handleSigterm executes when termination from operating system is received.
func (node *node) handleSigterm() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)

	<-sigc

	log.Warn().Msg("node: shutting down")

	log.Fatal().Msg("node: terminated")
}
