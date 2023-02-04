package main

import (
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"backend/blockchain"
	"backend/errors"
	"backend/networking"
	"backend/util"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/rs/zerolog/log"
)

// version is interpolated during build time.
var version string

// blocks stores all received blockchain blocks from
// other nodes within the network.
var blocks = make([][]blockchain.Block, 0)

// Node represents a singular blockchain node.
type Node struct {
	Version    string
	Uptime     time.Time
	interval   time.Duration
	network    *networking.Network
	blockchain *blockchain.Blockchain
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

	interval, err := time.ParseDuration(config.Interval)
	if err != nil {
		return nil, err
	}

	return &Node{
		Version:    version,
		interval:   interval,
		network:    net,
		blockchain: blockchain.NewBlockchain(),
		ready:      make(chan struct{}),
		close:      make(chan struct{}),
	}, nil
}

// Run starts all services required by the Node.
func (n *Node) Run() {
	// start network
	if err := n.network.Start(); err != nil {
		log.Fatal().Err(err).Msg("node: network failed to run")
	}

	// setup network stream handlers
	n.setStreamHandlers()

	// setup and initialize the blockchain
	n.setup()

	// wait for ready signal from setup
	<-n.ready

	// start listener for incoming requests
	n.listen()

	// start the scheduler
	if err := n.schedule(n.interval); err != nil {
		log.Fatal().Err(err).Msg("node: failed to start scheduler")
	}

	// node done provisioning; set uptime for node
	n.Uptime = time.Now()
}

// Stop tries to stop all running services of the Node.
// The network will be gracefully closed, and the current ledger of the blockchain
// will be written to the host.
func (n *Node) Stop() {
	close(n.close)

	if err := n.network.Close(); err != nil {
		log.Error().Err(err).Msg("node: failed to close network")
	}

	n.network.Host.RemoveStreamHandler("/reply")

	if err := n.blockchain.DumpJSON(); err != nil {
		log.Error().Err(err).Msg("node: failed to dump blockchain")
	}

	n.wg.Wait()
}

// setStreamHandlers sets the stream handlers that will handle individual request from other nodes.
func (n *Node) setStreamHandlers() {
	n.network.Host.SetStreamHandler("/reply", func(s network.Stream) {
		var message networking.Message

		msg, err := io.ReadAll(s)
		if err != nil {
			log.Error().Err(err).Msg("network: failed to read reply")
		}

		util.JSONDecode(msg, &message)

		log.Debug().Str("topic", string(message.Topic)).Msg("network: received reply")

		switch message.Topic {
		case networking.Blockchain:
			var b blockchain.Blockchain

			util.JSONDecode(message.Payload, &b)

			if len(b.Blocks) > 0 {
				blocks = append(blocks, b.Blocks)
			}
		case networking.Consensus:
			// TODO
		case networking.Block, networking.Transaction:
		}
	})
}

// setup will set up the blockchain from either scratch or by using a file containing
// the last known blockchain blocks from this current node (e.g. the Node has been shutdown,
// and is in the process of being rebooted). The file will only be used as a reference; all
// nodes within the network will be asked to send their current copy of the ledger. The file
// will only be used as the actual blockchain if the Node is either not connected to other nodes,
// or that the other nodes' blockchain is invalid; e.g. the integrity cannot be verified.
// The Node will be blocked from execution until a signal is given that the setup has been
// successfully completed.
func (n *Node) setup() {
	log.Info().
		Int("node(s)", n.network.ConnectedPeers()).
		Msg("node: synchronizing")

	// get blocks from file
	if data, err := n.blockchain.FromFile(); err == nil {
		if len(data) != 0 {
			blocks = append(blocks, data)
		}
	}

	// get blocks from peers
	time.AfterFunc(time.Second, func() {
		if n.network.ConnectedPeers() > 0 {
			n.network.Request(networking.Blockchain)
		}
	})

	n.wg.Add(1)

	// wait for replies
	time.AfterFunc(5*time.Second, func() {
		defer n.wg.Done()

		b := make([]blockchain.Block, 0)

		if len(blocks) > 0 {
			for _, data := range blocks {
				if len(data) > len(b) {
					b = data
				}
			}
		}

		// initialize blockchain
		n.blockchain.Init(n.network.ID(), b)

		// reset blocks
		blocks = nil
	})

	n.wg.Wait()

	log.Info().
		Int("node(s)", n.network.ConnectedPeers()).
		Msg("node: synchronized")

	close(n.ready)
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

	log.Debug().Msg("node: scheduler started")

	return nil
}

// reply sends a reply to another node using the network's Reply method.
func (n *Node) reply(peer string, topic networking.Topic, payload []byte) {
	n.network.Reply(peer, topic, payload)
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
			case msg := <-net.Subs[networking.Transaction].Messages:
				var t blockchain.Transaction

				util.JSONDecode(msg.Payload, &t)

				n.blockchain.AddTransaction(t)
			case msg := <-net.Subs[networking.Block].Messages:
				var b blockchain.Block

				util.JSONDecode(msg.Payload, &b)

				n.blockchain.AddBlock(b, msg.Peer)
			case msg := <-net.Subs[networking.Blockchain].Messages:
				if len(n.blockchain.Blocks) > 0 {
					n.reply(msg.Peer, networking.Blockchain, util.JSONEncode(n.blockchain))
				}
			}
		}
	}()

	log.Debug().Msg("node: listener started")
}

// HandleSigterm executes when termination from operating system is received.
// An attempt to gracefully shut down all required services of the node will be made.
// Execution will block until signal has been received; executing thread will wait on channel.
func (n *Node) HandleSigterm() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)

	<-sigc

	log.Info().Msg("node: shutting down")

	n.Stop()

	log.Info().Msg("node: terminated")
}
