package blockchain

import (
	"crypto/sha256"

	"backend/errors"
)

// source: https://github.com/tensor-programming/golang-blockchain/tree/part_10

// tree represents a simple implementation of a Merkle tree.
type tree struct {
	root   *node
	leaves []*node
}

// node represents a singular node within the Merkle tree.
type node struct {
	parent *node
	left   *node
	right  *node
	hash   []byte
}

// newMerkleNode creates a new node to be used in the Merkle tree.
func newMerkleNode(left, right *node, data []byte) *node {
	n := node{}

	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		n.hash = hash[:]
	} else {
		prevHashes := append(left.hash, right.hash...) //nolint
		hash := sha256.Sum256(prevHashes)
		n.hash = hash[:]
	}

	n.left = left
	n.right = right

	return &n
}

// buildRoot creates the root of the Merkle tree.
func (t *tree) buildRoot() *node {
	nodes := t.leaves

	for len(nodes) > 1 {
		var parents []*node

		if len(nodes)%2 != 0 {
			nodes = append(nodes, nodes[len(nodes)-1])
		}

		for i := 0; i < len(nodes); i += 2 {
			n := newMerkleNode(nodes[i], nodes[i+1], nil)
			parents = append(parents, n)
			nodes[i].parent, nodes[i+1].parent = n, n
		}

		nodes = parents
	}

	return nodes[0]
}

// newMerkleTree creates a new Merkle tree.
// The Merkle tree is used for validating a set of data by using their hash.
// If a data entry were to change, the result would cascade up to the root, and thus
// the data would be invalidated.
func newMerkleTree(data [][]byte) (*tree, error) {
	t := &tree{
		leaves: make([]*node, 0, len(data)),
	}

	for _, h := range data {
		n := newMerkleNode(nil, nil, h)
		t.leaves = append(t.leaves, n)
	}

	if len(t.leaves) == 0 {
		return nil, errors.ErrInvalidInput("no nodes could be created from specified input")
	}

	t.root = t.buildRoot()

	return t, nil
}
