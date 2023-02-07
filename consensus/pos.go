package consensus

import (
	"math/rand"
	"sync"
	"time"

	"backend/blockchain"
	"backend/errors"
)

// This is a very simple implementation of proof of stake. Should probably be refactored,
// but due to time constraints this will not happen.

// ProofOfStake is the consensus algorithm used by the node.
// Probably should refactor this to use another struct, instead of putting it all in this struct.
type ProofOfStake struct {
	Transactions map[string]any // remove (and insert) from map is O(1) whilst removing from an array is O(n) (iterate through array).
	Validators   map[string]any
	Responses    []Resp
	stakers      map[string]blockchain.Coin
	sync.RWMutex
}

// Resp the response of the consensus.
// Not optimal, but it works (for now).
type Resp struct {
	Data  []byte
	Valid bool
}

// NewPoS creates a new proof of stake consensus instance.
func NewPoS() *ProofOfStake {
	return &ProofOfStake{
		stakers:      make(map[string]blockchain.Coin),
		Transactions: make(map[string]any),
		Validators:   make(map[string]any),
		Responses:    make([]Resp, 0),
	}
}

// Winner returns the validator that is allowed to forge a new block.
func (pos *ProofOfStake) Winner() (string, error) {
	var node string

	pos.Lock()
	defer pos.Unlock()

	pool := make([]string, 0, len(pos.stakers))

	for k, v := range pos.stakers {
		if v.Float64() > 0 {
			pool = append(pool, k)
		}
	}

	if len(pool) == 0 {
		return node, errors.ErrInvalidOperation("no stakers")
	}

	r := rand.New(rand.NewSource(time.Now().Unix()))
	node = pool[r.Intn(len(pool))]

	return node, nil
}

// GetStake returns the stake of a given node.
func (pos *ProofOfStake) GetStake(node string) (blockchain.Coin, error) {
	if !pos.Exists(node) {
		return blockchain.Coin{}, errors.ErrInvalidOperation("node does not exist")
	}

	return pos.stakers[node], nil
}

// Set sets the stake of a given node.
func (pos *ProofOfStake) Set(node string, stake float64) {
	pos.Lock()
	defer pos.Unlock()

	pos.stakers[node] = blockchain.ToCoin(stake)
}

// Update updates the stake of a given node.
func (pos *ProofOfStake) Update(node string, stake float64) error {
	if !pos.Exists(node) {
		return errors.ErrInvalidOperation("node does not exist")
	}

	pos.Lock()
	defer pos.Unlock()

	if 0 > pos.stakers[node].Float64()+stake {
		return errors.ErrInvalidOperation("stake cannot be negative")
	}

	pos.stakers[node] = pos.stakers[node].Add(stake)

	return nil
}

// Add adds a node.
func (pos *ProofOfStake) Add(node string, stake float64) error {
	if pos.Exists(node) {
		return errors.ErrInvalidOperation("node already exists")
	}

	if 0 > stake {
		return errors.ErrInvalidOperation("stake cannot be negative")
	}

	pos.Set(node, stake)

	return nil
}

// Clear clears the proof of stake.
func (pos *ProofOfStake) Clear() {
	pos.Lock()
	defer pos.Unlock()

	pos.stakers = make(map[string]blockchain.Coin)
}

// Exists checks if the node exists.
func (pos *ProofOfStake) Exists(node string) bool {
	_, ok := pos.stakers[node]

	return ok
}
