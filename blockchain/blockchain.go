package blockchain

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"time"

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
func (b *Blockchain) Init(blocks []Block) {
	if len(blocks) > 0 {
		b.Blocks = blocks
	} else {
		if err := b.createGenesis(); err != nil {
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
	if err := b.validate(block, validator); err != nil {
		log.Error().Err(err).Msg("blockchain: failed to add block")

		return
	}

	b.Blocks = append(b.Blocks, block)
}

// CreateBlock creates a new block.
func (b *Blockchain) CreateBlock(validator string, amount uint16) (Block, error) {
	transactions := b.mp.retrieve(amount)

	block, err := createBlock(validator, b.Blocks[len(b.Blocks)-1].hash(), transactions)
	if err != nil {
		log.Error().Err(err).Msg("blockchain: failed to create block")

		return Block{}, err
	}

	if err = b.mp.delete(transactions...); err != nil {
		log.Error().Err(err).Msg("blockchain: failed to delete transactions")

		return Block{}, nil
	}

	return block, nil
}

// AddTransaction adds a new transaction to the memory pool.
func (b *Blockchain) AddTransaction(transaction Transaction) {
	if err := b.mp.add(transaction); err != nil {
		log.Error().Err(err).Msg("blockchain: failed to add transaction")

		return
	}

	// update the account of the sender
	if err := b.am.update(transaction.Sender, -transaction.Amount); err != nil {
		log.Error().Err(err).Msg("blockchain: failed to update account model")

		return
	}

	// update or add the account of the receiver
	if err := b.am.update(transaction.Receiver, transaction.Amount); err != nil {
		log.Error().Err(err).Msg("blockchain: failed to update account model")
	} else if err = b.am.add(transaction.Receiver, transaction.Amount, 0); err != nil {
		log.Error().Err(err).Msg("blockchain: failed to add account to account model")
	}
}

// CreateTransaction creates a new transaction
func (b *Blockchain) CreateTransaction(sender string, receiver string, signature []byte, amount float64) (Transaction, error) {
	// check if sender exists
	tx, err := b.am.get(sender)
	if err != nil {
		return Transaction{}, err
	}

	// check if sender has sufficient funds
	if amount > tx.Balance.Float64() {
		return Transaction{}, fmt.Errorf("%w: insufficient funds", errInvalidTransaction)
	}

	//check if signature is valid FIXME
	//if !ecdsa.VerifyASN1(key, []byte(fmt.Sprintf("%s%s%f", sender, receiver, amount)), signature) {
	//	return Transaction{}, fmt.Errorf("%w: invalid signature", errInvalidTransaction)
	//}

	// create transaction
	t := Transaction{
		Sender:    sender,
		Receiver:  receiver,
		Signature: hex.EncodeToString(signature),
		Amount:    amount,
		Nonce:     tx.Transactions,
		Timestamp: time.Now().Unix(),
	}

	if b.mp.exists(t.string()) {
		return Transaction{}, fmt.Errorf("%w: duplicate transaction", errInvalidTransaction)
	}

	// update the account of the sender
	if err = b.am.update(sender, -amount); err != nil {
		return Transaction{}, err
	}

	// update or add the account of the receiver
	if b.am.exists(receiver) {
		if err = b.am.update(receiver, amount); err != nil {
			return Transaction{}, err
		}
	} else {
		if err = b.am.add(receiver, amount, 0); err != nil {
			return Transaction{}, err
		}
	}

	// add transaction to memory pool
	if err = b.mp.add(t); err != nil {
		return Transaction{}, err
	}

	return t, nil
}

// validate validates a singular block.
func (b *Blockchain) validate(block Block, validator string) error {
	last := b.Blocks[len(b.Blocks)-1]

	// compare hashes
	if hex.EncodeToString(last.hash()) != block.PrevHash {
		return fmt.Errorf("%w, %s", errInvalidBlock, "hash does not match")
	}

	// check timstamp
	if last.Timestamp > block.Timestamp {
		return fmt.Errorf("%w, %s", errInvalidBlock, "invalid timestamp")
	}

	// compare validator
	if block.Validator != validator {
		return fmt.Errorf("%w, %s", errInvalidBlock, "invalid validator")
	}

	return nil
}

// createGenesis creates the genesis block.
func (b *Blockchain) createGenesis() error {
	log.Debug().Msg("blockchain: creating genesis block")

	t := Transaction{
		Sender:    "",
		Receiver:  "genesis",
		Signature: "",
		Amount:    math.MaxUint64,
		Nonce:     0,
		Timestamp: time.Now().Unix(),
	}

	block, err := createBlock("genesis", []byte(""), []Transaction{t})
	if err != nil {
		return err
	}

	b.Blocks = append(b.Blocks, block)

	return nil
}

// FromFile returns all blocks that are written to the dumpfile.
func (b *Blockchain) FromFile() ([]Block, error) {
	var blockchain Blockchain

	data, err := os.ReadFile(dumpFile)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &blockchain)
	if err != nil {
		return nil, err
	}

	log.Debug().Msg("blockchain: reading from file")

	return blockchain.Blocks, nil
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

	if err = os.WriteFile(dumpFile, data, 0600); err != nil { //nolint
		return err
	}

	return nil
}
