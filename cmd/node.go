package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"backend/blockchain"
	"backend/consensus"
	"backend/errors"
	"backend/networking"
	"backend/util"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/rs/zerolog/log"
)

// This file is a mess; refactoring should be done.
// Today is not that day.

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
	pos        *consensus.ProofOfStake
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
		pos:        consensus.NewPoS(),
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

	// set initial stake
	n.pos.Set(n.network.ID(), 0)

	// node done provisioning; set uptime for node
	n.Uptime = time.Now()
}

// AddTransaction adds a new transaction to the memory pool.
func (n *Node) AddTransaction(transaction blockchain.Transaction) error {
	// check if sender exists
	tx, err := n.blockchain.GetAccount(transaction.Sender)
	if err != nil {
		log.Debug().Err(err).Msg("node: could not find account")

		return err
	}

	// check if sender has sufficient funds
	if transaction.Amount > tx.Balance.Float64() || 0 > transaction.Amount {
		log.Debug().Err(err).Msg("node: account has insufficient funds")

		return fmt.Errorf("%w: insufficient funds", blockchain.ErrInvalidTransaction)
	}

	// FIXME
	// validate signature
	if err = transaction.Verify(); err != nil {
		return err
	}

	// check if sender has sufficient funds
	if transaction.Amount > tx.Balance.Float64() || 0 > transaction.Amount {
		return fmt.Errorf("%w: insufficient funds", blockchain.ErrInvalidTransaction)
	}

	// update the memory pool
	if err = n.blockchain.UpdateMempool(transaction); err != nil {
		log.Debug().Err(err).Msg("node: could not add transaction to mempool")

		return err
	}

	// update sender
	if err = n.blockchain.UpdateAccountModel(transaction.Sender, -transaction.Amount); err != nil {
		log.Debug().Err(err).Msg("node: could not update account")

		return err
	}

	// TODO change stake transaction

	log.Debug().Msg("node: added transaction")

	return nil
}

// CreateTransaction creates a new Transaction.
func (n *Node) CreateTransaction(sender string, receiver string, signature []byte, amount float64, txType blockchain.TxType) (blockchain.Transaction, error) {
	// check if sender exists
	tx, err := n.blockchain.GetAccount(sender)
	if err != nil {
		log.Debug().Err(err).Msg("node: could not find account")

		return blockchain.Transaction{}, err
	}

	// check if sender has sufficient funds
	if amount > tx.Balance.Float64() || 0 > amount {
		log.Debug().Err(err).Msg("node: account has insufficient funds")

		return blockchain.Transaction{}, fmt.Errorf("%w: insufficient funds", blockchain.ErrInvalidTransaction)
	}

	// create transaction
	t := blockchain.Transaction{
		Sender:    sender,
		Receiver:  receiver,
		Signature: util.HexEncode(signature),
		Amount:    blockchain.ToCoin(amount).Float64(),
		Nonce:     tx.Transactions,
		Timestamp: time.Now().Unix(),
		Type:      txType,
	}

	// FIXME (use priv key in api as temp fix)
	// validate signature
	if err = t.Verify(); err != nil {
		log.Debug().Err(err).Msg("node: could not verify transaction")

		return blockchain.Transaction{}, err
	}

	// update the memory pool
	if err = n.blockchain.UpdateMempool(t); err != nil {
		log.Debug().Err(err).Msg("node: could not add transaction to mempool")

		return blockchain.Transaction{}, err
	}

	// set stake
	if t.Type == blockchain.Stake {
		n.pos.Transactions[t.String()] = struct{}{}

		if err = n.pos.Update(n.network.ID(), t.Amount); err != nil {
			log.Debug().Err(err).Msg("node: could not update stake")

			return blockchain.Transaction{}, err
		}
	}

	// update sender
	if err = n.blockchain.UpdateAccountModel(t.Sender, -t.Amount); err != nil {
		log.Debug().Err(err).Msg("node: could not update account")

		return blockchain.Transaction{}, err
	}

	// publish message
	n.network.Publish(networking.Transaction, util.JSONEncode(t))

	log.Debug().Msg("node: created transaction")

	return t, nil
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
			var r consensus.Resp

			util.JSONDecode(message.Payload, &r)

			n.pos.Responses = append(n.pos.Responses, r)
		case networking.Stake:
			if f, err := strconv.ParseFloat(string(message.Payload), 64); err == nil {
				n.pos.Set(message.Peer, f)
			}
		case networking.Block, networking.Transaction, networking.Validator:
			// ignore; requests are handled by the listener
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

		ticker := time.NewTicker(30 * time.Second) // FIXME do this to interval

		for {
			select {
			case <-ticker.C:
				// request stake from other nodes
				n.network.Request(networking.Stake)

				// wait for replies
				time.AfterFunc(5*time.Second, func() {
					validator, err := n.pos.Winner()
					if err != nil {
						// no stakers; new block will be created by this node
						validator = n.network.ID()
					}

					// publish the node that will create the block
					n.network.Publish(networking.Validator, []byte(validator))

					// if this node is the validator; create block
					if validator == n.network.ID() {
						n.forge()
					}
				})
			case <-n.close:
				ticker.Stop()

				return
			}
		}
	}()

	log.Debug().Msg("node: scheduler started")

	return nil
}

// forge Forges a new block.
func (n *Node) forge() {
	// wait a bit before forging block
	time.AfterFunc(5*time.Second, func() {
		// create block with a max of 1000 transactions, returns an error if there are no transactions
		block, err := n.blockchain.CreateBlock(n.network.ID(), 1000)
		if err != nil {
			log.Debug().Err(err).Msg("node: failed to create block")

			return
		}

		n.network.Publish(networking.Consensus, util.JSONEncode(block))

		// wait for (consensus) replies
		time.AfterFunc(5*time.Second, func() {
			responses := len(n.pos.Responses)
			valid := 0

			for _, v := range n.pos.Responses {
				if bytes.Equal(block.Hash(), v.Data) {
					if v.Valid {
						valid++
					}
				}
			}

			var value = 100

			if responses != 0 {
				value = valid / responses * 100
			}

			// hardcoded value; meaning that it will not pass if there are only two nodes
			// should be done differently
			if value >= 66 {
				n.blockchain.AddBlock(block, n.network.ID())

				for _, t := range block.Transactions {
					if _, ok := n.pos.Transactions[t.String()]; ok {
						if err = n.pos.Update(n.network.ID(), -t.Amount); err != nil {
							log.Debug().Err(err).Msg("node: failed to update stake")
						}
						delete(n.pos.Transactions, t.String())
					}
				}

				n.network.Publish(networking.Block, util.JSONEncode(block))
			}

			// reset stakers
			v, _ := n.pos.GetStake(n.network.ID())
			n.pos.Clear()
			n.pos.Set(n.network.ID(), v.Float64())

			// reset responses
			n.pos.Responses = make([]consensus.Resp, 0)

			// remove validator
			delete(n.pos.Validators, n.network.ID())
		})
	})
}

// reply sends a reply to another node using the network's Reply method.
func (n *Node) reply(peer string, topic networking.Topic, payload []byte) {
	n.network.Reply(peer, topic, payload)
}

// listen listens to incoming traffic from all nodes that this Node is connected to.
// note: refactor this.
func (n *Node) listen() {
	n.wg.Add(1)

	net := n.network

	go func() {
		defer n.wg.Done()

		for {
			select {
			case <-n.close:
				return
			case msg := <-net.Subs[networking.Transaction].Messages: // transaction
				var t blockchain.Transaction

				util.JSONDecode(msg.Payload, &t)

				if err := n.AddTransaction(t); err != nil {
					log.Error().Err(err).Msg("node: failed to add transaction")
				}
			case msg := <-net.Subs[networking.Block].Messages: // block
				var b blockchain.Block

				util.JSONDecode(msg.Payload, &b)

				if _, ok := n.pos.Validators[msg.Peer]; ok {
					delete(n.pos.Validators, msg.Peer)

					n.blockchain.AddBlock(b, msg.Peer)
				}
			case msg := <-net.Subs[networking.Blockchain].Messages: // blockchain
				if len(n.blockchain.Blocks) > 0 {
					n.reply(msg.Peer, networking.Blockchain, util.JSONEncode(n.blockchain))
				}
			case msg := <-net.Subs[networking.Stake].Messages: // stake
				if stk, err := n.pos.GetStake(n.network.ID()); err == nil {
					n.network.Reply(msg.Peer, networking.Stake, util.JSONEncode(stk.Float64()))
				}
			case msg := <-net.Subs[networking.Consensus].Messages: // consensus
				var b blockchain.Block

				util.JSONDecode(msg.Payload, &b)

				resp := &consensus.Resp{
					Data:  b.Hash(),
					Valid: false,
				}

				err := b.Validate(n.blockchain.Blocks[len(n.blockchain.Blocks)-1], msg.Peer)
				if err != nil {
					n.network.Reply(msg.Peer, networking.Consensus, util.JSONEncode(resp))

					return
				}

				resp.Valid = true

				n.network.Reply(msg.Peer, networking.Consensus, util.JSONEncode(resp))
			case msg := <-net.Subs[networking.Validator].Messages: // validator
				// append validator to array, to keep track of validators.
				n.pos.Validators[string(msg.Payload)] = struct{}{}

				// if this node is the validator; create block
				if string(msg.Payload) == n.network.ID() {
					n.forge()
				}
			}
		}
	}()

	log.Debug().Msg("node: listener started")
}

// handleSigterm executes when termination from operating system is received.
// An attempt to gracefully shut down all required services of the node will be made.
// Execution will block until signal has been received; executing thread will wait on channel.
func (n *Node) handleSigterm() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)

	<-sigc

	log.Info().Msg("node: shutting down")

	n.Stop()

	log.Info().Msg("node: terminated")
}
