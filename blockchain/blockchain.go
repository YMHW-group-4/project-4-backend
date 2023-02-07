package blockchain

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"time"

	"backend/crypto"
	"backend/util"

	"github.com/rs/zerolog/log"
)

// dumpFile the file name to whom the Blockchain should be written.
const dumpFile string = "blockchain.json"

// Blockchain holds all the blocks in the Blockchain.
type Blockchain struct {
	Blocks []Block `json:"blocks"`
	mp     *mempool
	am     *accountModel
}

// NewBlockchain creates a new Blockchain.
func NewBlockchain() *Blockchain {
	return &Blockchain{
		Blocks: make([]Block, 0),
		am:     newAccountModel(),
		mp:     newMempool(),
	}
}

// Init initializes the blockchain and its account model.
func (b *Blockchain) Init(validator string, blocks []Block) {
	if len(blocks) > 0 {
		b.Blocks = blocks
	} else {
		if err := b.createGenesis(validator); err != nil {
			log.Fatal().Err(err).Msg("blockchain: failed to create genesis")
		}
	}

	log.Debug().Msg("blockchain: initializing account model")

	if len(b.Blocks) > 0 {
		b.am.fromBlocks(b.Blocks...)
	}
}

// AddBlock adds a new block to the blockchain.
func (b *Blockchain) AddBlock(block Block, validator string) {
	if err := block.Validate(b.Blocks[len(b.Blocks)-1], validator); err != nil {
		log.Error().Err(err).Msg("blockchain: block is invalid")

		return
	}

	for _, t := range block.Transactions {
		if t.Type == Stake {
			if err := b.UpdateAccountModel(t.Sender, t.Amount); err != nil {
				log.Debug().Err(err).Msg("failed to update account")
			}
		} else {
			if err := b.UpdateAccountModel(t.Receiver, t.Amount); err != nil {
				log.Debug().Err(err).Msg("failed to update account")
			}
		}
	}

	if err := b.mp.delete(block.Transactions...); err != nil {
		log.Debug().Err(err).Msg("failed to remove transactions")
	}

	log.Info().Str("validator", validator).Msg("blockchain: added new block")

	b.Blocks = append(b.Blocks, block)
}

// CreateBlock creates a new block.
func (b *Blockchain) CreateBlock(validator string, amount uint32) (Block, error) {
	transactions := b.mp.retrieve(amount)

	block, err := newBlock(validator, b.Blocks[len(b.Blocks)-1].Hash(), transactions)
	if err != nil {
		return Block{}, err
	}

	return block, nil
}

// createGenesis creates the genesis block.
func (b *Blockchain) createGenesis(validator string) error {
	priv, pub, err := crypto.Genesis()
	if err != nil {
		return err
	}

	sign, err := crypto.Sign(priv, []byte("genesis"))
	if err != nil {
		return err
	}

	t := Transaction{
		Sender:    util.HexEncode(crypto.EncodePublicKey(pub)),
		Receiver:  util.HexEncode(crypto.EncodePublicKey(pub)),
		Signature: util.HexEncode(sign),
		Amount:    ToCoin(math.MaxUint64).Float64(),
		Nonce:     0,
		Timestamp: time.Now().Unix(),
		Type:      Exchange,
	}

	block, err := newBlock(validator, []byte(""), []Transaction{t})
	if err != nil {
		return err
	}

	b.Blocks = append(b.Blocks, block)

	log.Debug().Msg("blockchain: created genesis block")

	return nil
}

// FromFile returns all blocks that are written to the dumpfile.
func (b *Blockchain) FromFile() ([]Block, error) {
	var blockchain Blockchain

	data, err := os.ReadFile(dumpFile)
	if err != nil {
		return nil, err
	}

	util.JSONDecode(data, &blockchain)

	if len(blockchain.Blocks) == 0 {
		return nil, err
	}

	log.Debug().Msg("blockchain: reading from file")

	return blockchain.Blocks, nil
}

// UpdateMempool tries to update or add to the memory pool.
func (b *Blockchain) UpdateMempool(transaction Transaction) error {
	if b.mp.exists(transaction.String()) {
		return fmt.Errorf("%w: duplicate transaction", ErrInvalidTransaction)
	}

	if err := b.mp.add(transaction); err != nil {
		return err
	}

	return nil
}

// UpdateAccountModel tries to update or add to the account model.
func (b *Blockchain) UpdateAccountModel(key string, amount float64) error {
	if b.am.exists(key) {
		if err := b.am.update(key, amount); err != nil {
			return err
		}
	} else {
		if err := b.am.add(key, amount); err != nil {
			return err
		}
	}

	return nil
}

// GetAccount returns the account associated with the given key.
func (b *Blockchain) GetAccount(key string) (*Account, error) {
	return b.am.get(key)
}

// DumpJSON writes the current Blockchain to a JSON file.
func (b *Blockchain) DumpJSON() error {
	log.Debug().Msg("blockchain: writing to file")

	data, err := json.MarshalIndent(b, "", " ")
	if err != nil {
		return err
	}

	if err = os.WriteFile(dumpFile, data, 0o600); err != nil {
		return err
	}

	return nil
}
