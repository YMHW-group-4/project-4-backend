package blockchain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
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
		mp:     newMempool(),
		am:     newAccountModel(),
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
func (b *Blockchain) AddBlock(block Block) error {
	if err := b.validate(block); err != nil {
		return err
	}

	b.Blocks = append(b.Blocks, block)

	return nil
}

// AddTransaction adds a new transaction to the memory pool.
func (b *Blockchain) AddTransaction(transaction Transaction) error {
	return b.mp.add(transaction)
}

// CreateTransaction creates a new transaction
func (b *Blockchain) CreateTransaction(sender string, receiver string, signature []byte, amount float32) (Transaction, error) {
	// check if sender exists
	tx, err := b.am.get(sender)
	if err != nil {
		return Transaction{}, err
	}

	// check if sender has sufficient funds
	if amount > tx.balance {
		return Transaction{}, fmt.Errorf("%w: insufficient funds", errInvalidTransaction)
	}

	hash := []byte(fmt.Sprintf("%s%s%f", sender, receiver, amount))

	sigPublicKeyECDSA, _ := crypto.SigToPub(hash, signature)

	sigPublicKeyBytes := crypto.FromECDSAPub(sigPublicKeyECDSA)

	log.Debug().Msgf("gamming %v", sigPublicKeyBytes)

	//fmt.Printf("Recovered public key: %x\n", sigPublicKeyBytes)

	// derive key from signature

	//genPub, err := x509.ParsePKIXPublicKey([]byte(sender))
	//if err != nil {
	//	return Transaction{}, err
	//}

	//key := genPub.(*ecdsa.PublicKey)
	//key, err := crypto.SigToPub(hash, signature)
	//if err != nil {
	//	return Transaction{}, err
	//}

	//check if signature is valid
	//if !ecdsa.VerifyASN1(key, []byte(fmt.Sprintf("%s%s%f", sender, receiver, amount)), signature) {
	//	return Transaction{}, fmt.Errorf("%w: invalid signature", errInvalidTransaction)
	//}

	// create transaction
	t := Transaction{
		Sender:    sender,
		Receiver:  receiver,
		Signature: string(signature),
		Amount:    amount,
		Nonce:     tx.transactions,
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

// TODO
func (b *Blockchain) validate(block Block) error {
	last := b.Blocks[len(b.Blocks)-1]

	if res := bytes.Compare(last.hash(), []byte(block.PrevHash)); res != 0 {
		return fmt.Errorf("%w, %s", errInvalidBlock, "hash does not match")
	}

	if last.Timestamp > block.Timestamp {
		return fmt.Errorf("%w, %s", errInvalidBlock, "invalid timestamp")
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
		Amount:    math.MaxFloat32,
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

	return blockchain.Blocks, nil
}

// DumpJSON writes the current Blockchain to a JSON file.
func (b *Blockchain) DumpJSON() error {
	log.Debug().Msg("blockchain: writing dumpfile")

	data, err := json.MarshalIndent(b, "", " ")
	if err != nil {
		return err
	}

	if err = os.WriteFile(dumpFile, data, 0600); err != nil { //nolint
		return err
	}

	return nil
}
