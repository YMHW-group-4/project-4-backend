package util

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type EnvSuite struct {
	suite.Suite
}

func TestEnvSuite(t *testing.T) {
	suite.Run(t, &EnvSuite{})
}

func (s *EnvSuite) AfterTest(suiteName, testName string) {
	os.Unsetenv("PORT")
}

func (s *EnvSuite) TestGetEnv() {
	os.Setenv("PORT", "9090")

	val := GetEnv("PORT", 9091)
	s.Equal(9090, val)
}

func (s *EnvSuite) TestGetEnvFallback() {
	val := GetEnv("PORT", 9091)
	s.Equal(9091, val)
}
