package blockchain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var transactions = []Transaction{
	{
		Sender:    "mike",
		Receiver:  "bob",
		Signature: "signature",
		Amount:    100,
		Nonce:     1,
		Timestamp: 123456789,
	},
	{
		Sender:    "bob",
		Receiver:  "douglas",
		Signature: "signature",
		Amount:    250,
		Nonce:     1,
		Timestamp: 123456789,
	},
	{
		Sender:    "alice",
		Receiver:  "john",
		Signature: "signature",
		Amount:    100,
		Nonce:     1,
		Timestamp: 123456789,
	},
	{
		Sender:    "patrick",
		Receiver:  "steve",
		Signature: "signature",
		Amount:    1000,
		Nonce:     1,
		Timestamp: 123456789,
	},
}

func TestNewEmptyMerkleTree(t *testing.T) {
	_, err := newMerkleTree(nil)

	assert.NotNil(t, err)
}

//func TestNewMerkleTree(t *testing.T) {
//	tr, _ := newMerkleTree(transactions)
//
//	assert.NotNil(t, tr)
//}

//func TestEqualTreeRoots(t *testing.T) {
//	x := append([]hashable{}, transactions...)
//
//	data := hash([]{transactions})
//	data := hashTransactions(transactions)
//	tr1, _ := newMerkleTree(data)
//	tr2, _ := newMerkleTree(data)
//
//	assert.Equal(suite.T(), tr1.root.hash, tr2.root.hash)
//}
//
//func (suite *MerkleTestSuite) TestNonEqualTreeRoots() {
//
//	x := copy(cData, suite.data)
//
//	log.Debug().Msgf("%d", x)
//
//	//cData[len(cData)-1] = Transaction{
//	//	Sender:    "patrick",
//	//	Receiver:  "steve",
//	//	Signature: "signature",
//	//	Amount:    375, // 1000 -> 375
//	//	Nonce:     1,
//	//	Timestamp: 123456789,
//	//}
//
//	log.Debug().Msgf("%v", suite.data)
//	log.Debug().Msgf("%v", cData)
//
//	tr1, _ := newMerkleTree(suite.data)
//	tr2, _ := newMerkleTree(cData)
//
//	assert.NotEqual(suite.T(), tr1.root.hash, tr2.root.hash)
//}
//
//func (suite *MerkleTestSuite) TestEqualTreeRootsSerialized() {
//	var b Block
//
//	block, _ := CreateBlock(suite.data, []byte("none")) //nolint
//
//}

//
//	s, _ := json.Marshal(block) //nolint
//	_ = json.Unmarshal(s, &b)   //nolint
//
//	data := hashTransactions(transactions)
//	sData := hashTransactions(b.Transactions)
//
//	tr1, _ := newMerkleTree(data)
//	tr2, _ := newMerkleTree(sData)
//
//	assert.Equal(t, tr1.root.hash, tr2.root.hash)
//}
//
////func TestNewMerkleNode(t *testing.T) {
////	x := []hashable{
////		[]byte("node1"),
////	}
////
////	data := [][]byte{
////		[]byte("node1"),
////		[]byte("node2"),
////		[]byte("node3"),
////		[]byte("node4"),
////		[]byte("node5"),
////		[]byte("node6"),
////		[]byte("node7"),
////	}
////
////	// level 1
////	mn1 := newMerkleNode(nil, nil, data[0])
////	mn2 := newMerkleNode(nil, nil, data[1])
////	mn3 := newMerkleNode(nil, nil, data[2])
////	mn4 := newMerkleNode(nil, nil, data[3])
////	mn5 := newMerkleNode(nil, nil, data[4])
////	mn6 := newMerkleNode(nil, nil, data[5])
////	mn7 := newMerkleNode(nil, nil, data[6])
////	mn8 := newMerkleNode(nil, nil, data[6])
////
////	// level 2
////	mn9 := newMerkleNode(mn1, mn2, nil)
////	mn10 := newMerkleNode(mn3, mn4, nil)
////	mn11 := newMerkleNode(mn5, mn6, nil)
////	mn12 := newMerkleNode(mn7, mn8, nil)
////
////	// level 3
////	mn13 := newMerkleNode(mn9, mn10, nil)
////	mn14 := newMerkleNode(mn11, mn12, nil)
////
////	// level 4
////	mn15 := newMerkleNode(mn13, mn14, nil)
////
////	root := fmt.Sprintf("%x", mn15.hash)
////	tr, _ := newMerkleTree(suite.data)
////
////	assert.Equal(t, root, fmt.Sprintf("%x", tr.root.hash))
////}
