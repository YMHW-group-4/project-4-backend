package util

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
)

type EncodingTestSuite struct {
	suite.Suite
	json testStruct
}

func (suite *EncodingTestSuite) SetupTest() {
	suite.json = testStruct{Something: "something"}
}

func TestEncodingSuite(t *testing.T) {
	suite.Run(t, new(EncodingTestSuite))
}

type testStruct struct {
	Something string `json:"something"`
}

func (suite *EncodingTestSuite) TestJSONEncode() {
	data1, _ := json.Marshal(suite.json)
	data2 := JSONEncode(suite.json)

	suite.Equal(data1, data2)
}

func (suite *EncodingTestSuite) TestJSONDecode() {
	var t1, t2 testStruct

	data1, _ := json.Marshal(suite.json)
	data2 := JSONEncode(suite.json)

	_ = json.Unmarshal(data1, &t1)

	JSONDecode(data2, &t2)

	suite.Equal(t1, t2)
}

func (suite *EncodingTestSuite) TestHexEncode() {
	input := []byte("input")
	data1 := hex.EncodeToString(input)
	data2 := HexEncode(input)

	suite.Equal(data1, data2)
}

func (suite *EncodingTestSuite) TestHexDecode() {
	input := hex.EncodeToString([]byte("input"))
	data1, _ := hex.DecodeString(input)
	data2 := HexDecode(input)

	suite.Equal(data1, data2)
}
