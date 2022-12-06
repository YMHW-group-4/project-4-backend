package networking

import (
	"context"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/rs/zerolog/log"
)

// Network represents a peer-to-peer network.
type Network struct {
	ctx   context.Context
	host  host.Host
	peers peers
}

// NewNetwork creates a new Network with given port.
func NewNetwork(port int) (*Network, error) {
	host, err := libp2p.New(libp2p.ListenAddrs())
	if err != nil {
		return nil, err
	}

	return &Network{
		ctx:   context.Background(),
		host:  host,
		peers: peers{},
	}, nil
}

// Start starts the network.
func (network Network) Start() {
}

// Close closes the network.
func (network Network) Close() {
	if err := network.host.Close(); err != nil {
		log.Warn().Err(err).Msg("network: could not close")
	}
}
