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
