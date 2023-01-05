package blockchain

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

const dumpFile string = "blockchain.json"

type Blockchain struct {
	Blocks []Block `json:"blocks"`

	wg    sync.WaitGroup
	rw    sync.RWMutex // NOTE Lock might be required for concurrency, not sure yet.
	close chan struct{}
}

// NewBlockchain creates a new Blockchain.
func NewBlockchain() *Blockchain {
	blockchain := &Blockchain{
		Blocks: make([]Block, 0),
		wg:     sync.WaitGroup{},
		rw:     sync.RWMutex{},
		close:  make(chan struct{}, 0),
	}

	return blockchain
}

// Setup reads Blocks from a file and requests the Blockchain from other nodes.
func (blockchain *Blockchain) Setup() {
	var b Blockchain

	// read blockchain from file.
	if data, err := os.ReadFile(dumpFile); err == nil {
		if err = json.Unmarshal(data, &b); err == nil {
			// FIXME check the blockchain
			blockchain.Blocks = b.Blocks
		}
	}

	// TODO get from nodes

	// TODO create genesis block
}

// schedule starts an internal ticker with given interval.
func (blockchain *Blockchain) schedule(interval time.Duration) error {
	if interval == 0 {
		return ErrInvalidArgument("interval can only be non-zero")
	}

	blockchain.wg.Add(1)

	go func() {
		defer blockchain.wg.Done()

		ticker := time.NewTicker(interval)

		for {
			select {
			case <-ticker.C:
				// TODO create new Block
				// Proof of stake
			case <-blockchain.close:
				ticker.Stop()

				return
			}
		}
	}()

	return nil
}

// Close closes and stops all running processes of the Blockchain.
func (blockchain *Blockchain) Close() {
	close(blockchain.close)

	blockchain.dumpJSON()
	blockchain.wg.Wait()
}

// dumpJSON writes the current Blockchain to a JSON file.
func (blockchain *Blockchain) dumpJSON() {
	data, err := json.MarshalIndent(blockchain, "", " ")
	if err != nil {
		log.Error().Err(err).Msg("blockchain: failed to marshal")
	}

	if err = os.WriteFile(dumpFile, data, 0600); err != nil { //nolint
		log.Error().Err(err).Msg("blockchain: failed to write to file")
	}
}
