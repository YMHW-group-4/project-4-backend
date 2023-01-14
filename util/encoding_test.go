package util

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
)

type EncodingTestSuite struct {
	suite.Suite
	test testStruct
}

func (suite *EncodingTestSuite) SetupTest() {
	suite.test = testStruct{Something: "something"}
}

func TestEncodingSuite(t *testing.T) {
	suite.Run(t, new(EncodingTestSuite))
}

type testStruct struct {
	Something string `json:"something"`
}

func (suite *EncodingTestSuite) TestMarshalType() {
	data1, _ := json.Marshal(suite.test) //nolint
	data2 := MarshalType(suite.test)

	suite.Equal(data1, data2)
}

func (suite *EncodingTestSuite) TestUnmarshalType() {
	var t1, t2 testStruct

	data1, _ := json.Marshal(suite.test) //nolint
	data2 := MarshalType(suite.test)

	_ = json.Unmarshal(data1, &t1) //nolint

	UnmarshalType(data2, &t2)

	suite.Equal(t1, t2)
}
