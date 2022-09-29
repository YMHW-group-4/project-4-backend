package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGreet(t *testing.T) {
	assert.Equal(t, "Hello John!", greet("John"))
}

func TestGreetFail(t *testing.T) {
	assert.NotEqual(t, "Hello John!", greet("Steve"))
}
