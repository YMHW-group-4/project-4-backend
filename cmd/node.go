package main

import (
	"encoding/json"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"backend/api"
	"backend/blockchain"
	"backend/errors"
	"backend/networking"

	"github.com/libp2p/go-libp2p/core/network" // FIXME do not import this here; put logic in network.
	"github.com/rs/zerolog/log"
)

// version is interpolated during build time.
var version string

// Node represents a singular blockchain node.
type Node struct {
	Version    string
	Uptime     time.Time
	network    *networking.Network
	blockchain *blockchain.Blockchain
	api        *api.API
	rw         sync.RWMutex
	wg         sync.WaitGroup
	ready      chan struct{}
	close      chan struct{}
}

// NewNode creates a new Node with given configuration.
func NewNode(config Configuration) (*Node, error) {
	net, err := networking.NewNetwork(config.Port)
	if err != nil {
		return nil, err
	}

	return &Node{
		Version:    version,
		network:    net,
		blockchain: blockchain.NewBlockchain(),
		api:        api.NewAPI(config.ApiPort),
		ready:      make(chan struct{}, 1),
		close:      make(chan struct{}, 0),
	}, nil
}

// Run starts all services required by the Node.
func (n *Node) Run() {
	if err := n.network.Start(); err != nil {
		log.Fatal().Err(err).Msg("node: network failed to run")
	}

	if err := n.api.Start(); err != nil {
		log.Fatal().Err(err).Msg("node: api failed to run")
	}

	n.setup()

	<-n.ready

	if err := n.schedule(time.Minute * 10); err != nil {
		log.Fatal().Err(err).Msg("node: blockchain failed to run")
	}

	n.Uptime = time.Now()
}

// Stop tries to stop all running services of the Node.
// The network will be gracefully closed, and the current ledger of the blockchain
// will be writen to the host.
func (n *Node) Stop() {
	close(n.close)

	if err := n.network.Close(); err != nil {
		log.Error().Err(err).Msg("node: failed to close network")
	}

	if err := n.blockchain.DumpJSON(); err != nil {
		log.Error().Err(err).Msg("node: failed to dump blockchain")
	}

	n.wg.Wait()
}

// setup will set up the blockchain from either scratch or by using a file containing
// the last known blockchain blocks from this current node (e.g. the Node has been shutdown,
// and is in the process of being rebooted). The file will only be used as a reference; all
// nodes within the network will be asked to send their current copy of the ledger. The file
// will only be used as the actual blockchain if the Node is either not connected to other nodes,
// or that the other nodes' blockchain is invalid; e.g. the integrity cannot be verified.
// The Node will be blocked from execution until a signal is given that the setup has been
// succesfully completed.
func (n *Node) setup() {
	// create genesis block
	n.blockchain.CreateGenesis()

	// get blocks from file
	if _, err := n.blockchain.BlocksFromFile(); err == nil {
		// TODO
	}

	// get blocks from peers
	if n.network.ConnectedPeers() > 1 {
		if err := n.network.Request(networking.Blockchain); err == nil {
			// TODO wait for responses
		}
	}

	time.Sleep(5 * time.Second) // remove this

	n.ready <- struct{}{}
}

// schedule starts an internal ticker with given interval.
// each interval an attempt will be made to forge a new block on the blockchain.
func (n *Node) schedule(interval time.Duration) error {
	if interval == 0 {
		return errors.ErrInvalidArgument("interval can only be non-zero")
	}

	n.wg.Add(1)

	go func() {
		defer n.wg.Done()

		ticker := time.NewTicker(interval)

		for {
			select {
			case <-ticker.C:
				// TODO create new Block
				// Proof of stake
			case <-n.close:
				ticker.Stop()

				return
			}
		}
	}()

	return nil
}

// listen listens to incoming traffic from all nodes that this Node is connected to.
func (n *Node) listen() {
	n.wg.Add(1)

	net := n.network

	go func() {
		defer n.wg.Done()

		for {
			select {
			case <-n.close:
				return
			case _ = <-net.Subs[networking.Transaction].Messages:
				// do something
			case _ = <-net.Subs[networking.Block].Messages:
				// do something
			case _ = <-net.Subs[networking.Blockchain].Messages:
				// do something
			}
		}
	}()

	net.Host.SetStreamHandler("/reply", func(s network.Stream) {
		var message networking.Message

		b, err := io.ReadAll(s)
		if err != nil {
			log.Error().Err(err).Msg("network: failed to read reply")
		}

		if err = json.Unmarshal(b, &message); err != nil {
			log.Error().Err(err).Msg("network: failed to unmarshal reply")
		}

		switch message.Topic {
		case networking.Transaction:
			// do something
		case networking.Block:
			// do something
		case networking.Blockchain:
			// do something
		}
	})
}

// handleSigterm executes when termination from operating system is received.
// An attempt to gracefully shut down all required services of the node will be made.
func (n *Node) handleSigterm() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)

	<-sigc

	log.Info().Msg("node: shutting down")

	n.Stop()

	log.Info().Msg("node: terminated")
}
