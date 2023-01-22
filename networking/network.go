package networking

import (
	"context"
	"fmt"
	"sync"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
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
	host *Network
}

// Network represents a peer-to-peer network.
type Network struct {
	Host   host.Host
	Subs   map[Topic]*Subscription
	ctx    context.Context
	ps     *pubsub.PubSub
	wg     sync.WaitGroup
	close  chan struct{}
	notify chan struct{}
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
		Host:   h,
		Subs:   make(map[Topic]*Subscription, 0),
		ctx:    ctx,
		ps:     ps,
		wg:     sync.WaitGroup{},
		close:  make(chan struct{}),
		notify: make(chan struct{}),
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

	return nil
}

// ConnectedPeers returns the amount of currently connected peers.
func (n *Network) ConnectedPeers() int {
	return len(n.Host.Network().Peers())
}

// Publish publishes a Message to given Topic.
func (n *Network) Publish(topic Topic, payload []byte) error {
	msg, err := NewMessage(n.Host.ID().String(), topic, payload)
	if err != nil {
		return err
	}

	if err = n.Subs[topic].Publish(msg); err != nil {
		return err
	}

	return nil
}

// Request alias for Publish to quickly make a request on a Topic.
func (n *Network) Request(topic Topic) error {
	if err := n.Publish(topic, []byte("request")); err != nil {
		return err
	}

	return nil
}

// Reply sends a Message to one given peer.
func (n *Network) Reply(node string, topic Topic, payload []byte) error {
	id, err := peer.Decode(node)
	if err != nil {
		return err
	}

	s, err := n.Host.NewStream(n.ctx, id, "/reply")
	if err != nil {
		return err
	}

	msg, err := NewMessage(n.Host.ID().String(), topic, payload)
	if err != nil {
		return err
	}

	if _, err = s.Write(msg); err != nil {
		return err
	}

	if err = s.Close(); err != nil {
		return err
	}

	return nil
}

// Close closes the Network.
func (n *Network) Close() error {
	close(n.close)

	for _, sub := range n.Subs {
		if err := sub.Close(); err != nil {
			return err
		}
	}

	if err := n.Host.Close(); err != nil {
		return err
	}

	n.wg.Wait()

	return nil
}

// startMdns creates and starts a new mDNS service.
// This automatically discovers peers on the same LAN and connects to them.
func (n *Network) startMdns() error {
	s := mdns.NewMdnsService(n.Host, discoveryServiceTag, &discoveryNotifee{host: n})

	return s.Start()
}

// setupSubscriptions starts and listens to all Subscriptions.
func (n *Network) setupSubscriptions() error {
	for _, top := range []Topic{Transaction, Block, Blockchain} {
		sub, err := NewSubscription(n.ctx, n.ps, n.Host.ID(), top)
		if err != nil {
			return err
		}

		n.Subs[top] = sub
	}

	return nil
}

// Notify blocks until a peer has connected to the host.
func (n *Network) Notify() <-chan struct{} {
	return n.notify
}

// HandlePeerFound gets called when a new peer is discovered.
// This will automatically connect with the discovered peer.
func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	if err := n.host.Host.Connect(context.Background(), pi); err != nil {
		log.Error().Err(err).Msg("network: failed to connect to peer")
	}

	log.Debug().Msg("network: connected to peer")

	n.host.notify <- struct{}{}
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
