package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfigFromEnv(t *testing.T) {
	config := getConfigFromEnv()

	assert.Equal(t, false, config.debug)
	assert.Equal(t, 30333, config.port)
}
