package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfigFromEnv(t *testing.T) {
	config := getConfigFromEnv()

	assert.Equal(t, false, config.Debug)
	assert.Equal(t, 30333, config.Port)
	assert.Equal(t, 8080, config.APIPort)
	assert.Equal(t, "20m", config.Interval)
}
