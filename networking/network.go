package networking

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/routing"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/libp2p/go-libp2p/p2p/security/noise"
	"github.com/libp2p/go-libp2p/p2p/security/tls"
	libp2pquic "github.com/libp2p/go-libp2p/p2p/transport/quic"
	"github.com/rs/zerolog/log"
)

// discoveryServiceTag is used in mDNS advertisements to discover other peers.
const discoveryServiceTag = "crypto"

// discoveryNotifee gets notified when a new peer is discovered via mDNS.
type discoveryNotifee struct {
	host host.Host
}

// Network represents a peer-to-peer network.
type Network struct {
	ctx   context.Context
	host  host.Host
	ps    *pubsub.PubSub
	subs  map[Topic]*Subscription
	wg    sync.WaitGroup
	close chan struct{}
}

// NewNetwork creates a new Network with given port.
func NewNetwork(port int) (*Network, error) {
	ctx := context.Background()
	h, err := libp2p.New(
		libp2p.ListenAddrStrings(hostAddr(port)...),
		libp2p.Security(libp2ptls.ID, libp2ptls.New),
		libp2p.Security(noise.ID, noise.New),
		libp2p.Transport(libp2pquic.NewTransport),
		libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			idht, err := dht.New(ctx, h)

			return idht, err
		}),
	)
	if err != nil { //nolint
		return nil, err
	}

	ps, err := pubsub.NewGossipSub(ctx, h)
	if err != nil {
		return nil, err
	}

	return &Network{
		ctx:   ctx,
		host:  h,
		ps:    ps,
		subs:  make(map[Topic]*Subscription, 0),
		wg:    sync.WaitGroup{},
		close: make(chan struct{}, 0),
	}, nil
}

// Start starts the Network.
func (n *Network) Start() error {
	if err := n.setupSubscriptions(); err != nil {
		return err
	}

	if err := n.startMdns(); err != nil {
		return err
	}

	n.listen()

	return nil
}

// Publish publishes a Message to given Topic.
func (n *Network) Publish(topic Topic, payload string) {
	msg, err := NewMessage(n.host.ID().String(), topic, payload)
	if err != nil {
		log.Error().Err(err).Msg("network: failed to create message")
	}

	if err = n.subs[topic].Publish(msg); err != nil {
		log.Error().Err(err).Msg("network: failed to publish message")
	}
}

// Reply sends a Message to one given peer.
func (n *Network) Reply(peer peer.ID, topic Topic, payload string) {
	s, err := n.host.NewStream(n.ctx, peer, "/reply")
	if err != nil {
		log.Error().Err(err).
			Str("peer", peer.String()).
			Msg("network: failed to setup reply stream")
	}

	msg, err := NewMessage(n.host.ID().String(), topic, payload)
	if err != nil {
		log.Error().Err(err).
			Str("peer", peer.String()).
			Msg("network: failed to create reply message")
	}

	if _, err = s.Write(msg); err != nil {
		log.Error().Err(err).
			Str("peer", peer.String()).
			Msg("network: failed to write reply")
	}

	if err = s.Close(); err != nil {
		log.Error().Err(err).
			Str("peer", peer.String()).
			Msg("network: failed to close reply stream")
	}
}

// Close closes the Network.
func (n *Network) Close() error {
	close(n.close)
	n.host.RemoveStreamHandler("/reply")

	for _, sub := range n.subs {
		if err := sub.Close(); err != nil {
			return err
		}
	}

	if err := n.host.Close(); err != nil {
		return err
	}

	n.wg.Wait()

	return nil
}

// startMdns creates and starts a new mDNS service.
// This automatically discovers peers on the same LAN and connects to them.
func (n *Network) startMdns() error {
	s := mdns.NewMdnsService(n.host, discoveryServiceTag, &discoveryNotifee{host: n.host})

	return s.Start()
}

// setupSubscriptions starts and listens to all Subscriptions.
func (n *Network) setupSubscriptions() error {
	for _, top := range []Topic{Transaction, Block} {
		sub, err := NewSubscription(n.ctx, n.ps, n.host.ID(), top)
		if err != nil {
			return err
		}

		n.subs[top] = sub
	}

	return nil
}

// listen listens to incoming Messages from all Subscriptions and replies from Nodes.
func (n *Network) listen() {
	n.wg.Add(1)

	go func() {
		defer n.wg.Done()

		for {
			select {
			case <-n.close:
				return
			case msg := <-n.subs[Transaction].Messages:
				log.Debug().
					Str("topic", string(Transaction)).
					Str("payload", msg.Payload).
					Str("peer", msg.Peer).
					Msg("network: received message")
			case msg := <-n.subs[Block].Messages:
				log.Debug().
					Str("topic", string(Block)).
					Str("payload", msg.Payload).
					Str("peer", msg.Peer).
					Msg("network: received message")
			}
		}
	}()

	n.host.SetStreamHandler("/reply", func(s network.Stream) {
		var message Message

		b, err := io.ReadAll(s)
		if err != nil {
			log.Error().Err(err).Msg("network: failed to read reply")
		}

		if err = json.Unmarshal(b, &message); err != nil {
			log.Error().Err(err).Msg("network: failed to unmarshal reply")
		}

		switch message.Topic {
		case Transaction:
			// do something
		case Block:
			// do something
		}
	})
}

// HandlePeerFound gets called when a new peer is discovered.
// This will automatically connect with the discovered peer.
func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	if err := n.host.Connect(context.Background(), pi); err != nil {
		log.Debug().
			Err(err).
			Str("peer", pi.String()).
			Msg("network: failed to connect to peer")
	}

	log.Debug().
		Str("peer", pi.String()).
		Int("peer(s)", len(n.host.Network().Peers())).
		Msg("network: discovered peer")
}

// hostAddr makes address on given input ports for IPv4 and IPv6.
func hostAddr(ports ...int) []string {
	addrs := make([]string, 0, len(ports))

	for _, port := range ports {
		addrs = append(addrs, []string{
			fmt.Sprintf("/ip4/0.0.0.0/udp/%d/quic", port),
			fmt.Sprintf("/ip6/::/udp/%d/quic", port),
		}...)
	}

	return addrs
}
