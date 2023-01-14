package util

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type EnvTestSuite struct {
	suite.Suite
}

func TestEnvSuite(t *testing.T) {
	suite.Run(t, &EnvTestSuite{})
}

func (s *EnvTestSuite) AfterTest(suiteName, testName string) {
	_ = os.Unsetenv("NODE_PORT")
}

func (s *EnvTestSuite) TestGetEnv() {
	_ = os.Setenv("NODE_PORT", "9090")

	val := GetEnv("NODE_PORT", 9091)
	s.Equal(9090, val)
}

func (s *EnvTestSuite) TestGetEnvFallback() {
	val := GetEnv("NODE_PORT", 9091)
	s.Equal(9091, val)
}
