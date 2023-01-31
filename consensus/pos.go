package consensus

import "backend/blockchain"

type ProofOfStake struct {
	Stakers map[string]blockchain.Coin
}
