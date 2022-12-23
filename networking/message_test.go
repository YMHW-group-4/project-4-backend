package networking

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalJSON(t *testing.T) {
	var m Message

	message := Message{
		Peer:    "peer",
		Topic:   Transaction,
		Payload: "message",
	}

	msg, _ := NewMessage("peer", Transaction, "message")
	_ = json.Unmarshal(msg, &m)

	assert.Equal(t, message, m)
}
