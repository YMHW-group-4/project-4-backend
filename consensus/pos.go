package consensus

import "backend/blockchain"

type ProofOfStake struct {
	Stakers map[string]blockchain.Coin
}

// NewPoS creates a new proof of stake consensus instance.
func NewPoS() *ProofOfStake {
	return &ProofOfStake{
		Stakers: make(map[string]blockchain.Coin),
	}
}
