package networking

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
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
	ctx  context.Context
	host host.Host
	ps   *pubsub.PubSub
}

// NewNetwork creates a new Network with given port.
func NewNetwork(port int) (*Network, error) {
	ctx := context.Background()

	h, err := libp2p.New(libp2p.ListenAddrStrings(hostAddr(port)...))
	if err != nil {
		return nil, err
	}

	ps, err := pubsub.NewGossipSub(ctx, h)
	if err != nil {
		return nil, err
	}

	return &Network{
		ctx:  ctx,
		host: h,
		ps:   ps,
	}, nil
}

// Start starts the network.
func (network Network) Start() error {
	if err := network.setupMdns(); err != nil {
		return err
	}

	return nil
}

// Close closes the network.
func (network Network) Close() error {
	if err := network.host.Close(); err != nil {
		return err
	}

	return nil
}

// setupMdns creates and starts a new mDNS service.
// This automatically discovers peers on the same LAN and connects to them.
func (network Network) setupMdns() error {
	s := mdns.NewMdnsService(network.host, discoveryServiceTag, &discoveryNotifee{host: network.host})

	return s.Start()
}

// HandlePeerFound gets called when a new peer is discovered.
// This will automatically connect with the discovered peer.
func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	if err := n.host.Connect(context.Background(), pi); err != nil {
		log.Debug().Err(err).Msg("network: failed to connect")
	}

	log.Debug().
		Str("peer", pi.ID.String()).
		Msg("network: discovered peer")
}

// hostAddr makes address on given input ports for IPv4 and IPv6.
func hostAddr(ports ...int) []string {
	addrs := make([]string, 0, len(ports))

	for _, port := range ports {
		addrs = append(addrs, []string{
			fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port),
			fmt.Sprintf("/ip6/::/tcp/%d", port),
		}...)
	}

	return addrs
}
