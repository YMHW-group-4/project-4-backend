package blockchain

import "crypto/sha256"

// source: https://github.com/tensor-programming/golang-blockchain/tree/part_10

// tree represents a simple implementation of a Merkle tree.
type tree struct {
	root *node
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
		prevHashes := appendNodes(left, right)
		hash := sha256.Sum256(prevHashes)
		n.hash = hash[:]
	}

	n.left = left
	n.right = right

	return &n
}

// appendNodes returns the combined result of the hashes of two nodes.
func appendNodes(left, right *node) []byte {
	return append(left.hash, right.hash...)
}

// newMerkleTree creates a new Merkle tree.
// The Merkle tree is used for validating a set of data by using their hash.
// If a data entry were to change, the result would cascade up to the root, and thus
// the data would be invalidated.
func newMerkleTree(data [][]byte) (*tree, error) {
	nodes := make([]node, 0)

	for _, dat := range data {
		n := newMerkleNode(nil, nil, dat)
		nodes = append(nodes, *n)
	}

	if len(nodes) == 0 {
		return nil, ErrInvalidData("no nodes could be created from specified input")
	}

	for len(nodes) > 1 {
		var parents []node

		if len(nodes)%2 != 0 {
			nodes = append(nodes, nodes[len(nodes)-1])
		}

		for i := 0; i < len(nodes); i += 2 {
			n := newMerkleNode(&nodes[i], &nodes[i+1], nil)
			parents = append(parents, *n)
			nodes[i].parent, nodes[i+1].parent = n, n
		}

		nodes = parents
	}

	return &tree{&nodes[0]}, nil
}
