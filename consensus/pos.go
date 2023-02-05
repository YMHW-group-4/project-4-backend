package consensus

import (
	"backend/blockchain"
	"backend/errors"
)

// ProofOfStake is the consensus algorithm used by the node.
type ProofOfStake struct {
	stakers map[string]blockchain.Coin
}

// NewPoS creates a new proof of stake consensus instance.
func NewPoS() *ProofOfStake {
	return &ProofOfStake{
		stakers: make(map[string]blockchain.Coin),
	}
}

// get keys of node
// h.Peerstore().PrivKey(h.ID())
// h.Peerstore().PubKey(...)

// Winner returns the validator that is allowed to forge a new block.
func (pos *ProofOfStake) Winner() (string, error) {
	var node string

	if len(pos.stakers) == 0 {
		return node, errors.ErrInvalidOperation("no stakers")
	}

	pool := make([]string, 0, len(pos.stakers))

	for k, v := range pos.stakers {
		if v.Float64() > 0 {
			pool = append(pool, k)
		}
	}

	return node, nil
}

// GetStake returns the stake of a given node.
func (pos *ProofOfStake) GetStake(node string) (blockchain.Coin, error) {
	if err := pos.Exists(node); err != nil {
		return blockchain.Coin{}, err
	}

	return pos.stakers[node], nil
}

// Update updates the stake of a given node.
func (pos *ProofOfStake) Update(node string, stake float64) error {
	if err := pos.Exists(node); err != nil {
		return err
	}

	if 0 > pos.stakers[node].Float64()+stake {
		return errors.ErrInvalidOperation("stake cannot be negative")
	}

	pos.stakers[node] = pos.stakers[node].Add(stake)

	return nil
}

// Add adds a node.
func (pos *ProofOfStake) Add(node string, stake float64) error {
	if err := pos.Exists(node); err != nil {
		return err
	}

	if 0 > stake {
		return errors.ErrInvalidOperation("stake cannot be negative")
	}

	pos.stakers[node] = blockchain.ToCoin(stake)

	return nil
}

// Clear clears the proof of stake.
func (pos *ProofOfStake) Clear() {
	pos.stakers = make(map[string]blockchain.Coin)
}

// Exists checks if the node exists.
func (pos *ProofOfStake) Exists(node string) error {
	if _, ok := pos.stakers[node]; !ok {
		return errors.ErrInvalidOperation("node does not exist")
	}

	return nil
}
