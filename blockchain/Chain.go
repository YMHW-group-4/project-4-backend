package blockchain

import (
	"crypto"
	"crypto/rsa"
	"fmt"
	"math/big"
)

type Chain []Block

func createChain() Chain {
	firstTransaction := CreateTransaction(999, "genesis", "Wizard")
	prevBlock := CreateBlock("nil", firstTransaction)
	chain := make(Chain, 0)
	return append(chain, prevBlock)
}

func (chain Chain) GetLastBlock() Block {
	count := len(chain)
	return chain[count-1]
}

func (chain Chain) AddBlock(transaction Transaction, publicKey crypto.PublicKey, signature *rsa.PrivateKey) Chain {
	newBlock := CreateBlock(chain.GetLastBlock().Hash(), transaction)
	return append(chain, newBlock)
}

func (chain Chain) AddBlockV2(transaction Transaction, publicKey crypto.PublicKey, signature *rsa.PrivateKey) Chain {
	// TODO: Is valid, using the publicKey and the Signature
	isValid := true
	if isValid {
		newBlock := CreateBlock(chain.GetLastBlock().Hash(), transaction)
		solution := chain.mine(newBlock.nonce)
		if chain.verify(newBlock.nonce, solution) {
			return append(chain, newBlock)
		}
	}
	return chain
}

func (chain Chain) mine(nonce int64) int64 {
	solution := int64(1)
	fmt.Print("mining...")
	// TODO: This is a shitty proof of work method.
	for true {
		hex := big.NewInt(nonce + solution).Text(16)
		fmt.Print(hex)
		if hex[1:4] == "FFFF" {
			return solution
		}
		solution++
	}
	return int64(0)
}

func (chain Chain) verify(nonce int64, solution int64) bool {
	hex := big.NewInt(nonce + solution).Text(16)
	fmt.Print(hex)
	return hex[1:4] == "FFFF"
}
