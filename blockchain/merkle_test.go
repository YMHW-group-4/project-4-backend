package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var transactions = []mockTransaction{
	{From: "mike", To: "bob", Value: "100"},
	{From: "bob", To: "douglas", Value: "250"},
	{From: "alice", To: "john", Value: "100"},
	{From: "patrick", To: "steve", Value: "1000"},
}

type mockBlock struct {
	Transactions []mockTransaction `json:"transactions"`
}

type mockTransaction struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value string `json:"value"`
}

func hashTransactions(t []mockTransaction) [][]byte {
	data := make([][]byte, 0)
	h := sha256.New()

	for _, tx := range t {
		h.Write([]byte(fmt.Sprintf("%v", tx)))
		data = append(data, h.Sum(nil))
	}

	return data
}

func TestNewEmptyMerkleTree(t *testing.T) {
	_, err := newMerkleTree(nil)

	assert.NotNil(t, err)
}

func TestNewMerkleTree(t *testing.T) {
	data := hashTransactions(transactions)
	tr, _ := newMerkleTree(data)

	assert.NotNil(t, tr)
}

func TestEqualTreeRoots(t *testing.T) {
	data := hashTransactions(transactions)
	tr1, _ := newMerkleTree(data)
	tr2, _ := newMerkleTree(data)

	assert.Equal(t, tr1.root.hash, tr2.root.hash)
}

func TestNonEqualTreeRoots(t *testing.T) {
	cTransactions := []mockTransaction{
		{From: "mike", To: "bob", Value: "100"},
		{From: "bob", To: "douglas", Value: "250"},
		{From: "alice", To: "john", Value: "100"},
		{From: "patrick", To: "steve", Value: "375"}, // 1000 -> 375
	}

	data := hashTransactions(transactions)
	cData := hashTransactions(cTransactions)
	tr1, _ := newMerkleTree(data)
	tr2, _ := newMerkleTree(cData)

	assert.NotEqual(t, tr1.root.hash, tr2.root.hash)
}

func TestEqualTreeRootsSerialized(t *testing.T) {
	var b mockBlock

	block := &mockBlock{transactions}

	s, _ := json.Marshal(block) //nolint
	_ = json.Unmarshal(s, &b)   //nolint

	data := hashTransactions(transactions)
	sData := hashTransactions(b.Transactions)

	tr1, _ := newMerkleTree(data)
	tr2, _ := newMerkleTree(sData)

	assert.Equal(t, tr1.root.hash, tr2.root.hash)
}

func TestNewMerkleNode(t *testing.T) {
	data := [][]byte{
		[]byte("node1"),
		[]byte("node2"),
		[]byte("node3"),
		[]byte("node4"),
		[]byte("node5"),
		[]byte("node6"),
		[]byte("node7"),
	}

	// level 1
	mn1 := newMerkleNode(nil, nil, data[0])
	mn2 := newMerkleNode(nil, nil, data[1])
	mn3 := newMerkleNode(nil, nil, data[2])
	mn4 := newMerkleNode(nil, nil, data[3])
	mn5 := newMerkleNode(nil, nil, data[4])
	mn6 := newMerkleNode(nil, nil, data[5])
	mn7 := newMerkleNode(nil, nil, data[6])
	mn8 := newMerkleNode(nil, nil, data[6])

	// level 2
	mn9 := newMerkleNode(mn1, mn2, nil)
	mn10 := newMerkleNode(mn3, mn4, nil)
	mn11 := newMerkleNode(mn5, mn6, nil)
	mn12 := newMerkleNode(mn7, mn8, nil)

	// level 3
	mn13 := newMerkleNode(mn9, mn10, nil)
	mn14 := newMerkleNode(mn11, mn12, nil)

	// level 4
	mn15 := newMerkleNode(mn13, mn14, nil)

	root := fmt.Sprintf("%x", mn15.hash)
	tr, _ := newMerkleTree(data)

	assert.Equal(t, root, fmt.Sprintf("%x", tr.root.hash))
}
