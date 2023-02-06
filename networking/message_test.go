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
		Payload: []byte("message"),
	}

	msg, _ := json.Marshal(NewMessage("peer", Transaction, []byte("message")))
	_ = json.Unmarshal(msg, &m)

	assert.Equal(t, message, m)
}

func TestMarshalType(t *testing.T) {
	type testPassenger struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	type testBus struct {
		Passengers []testPassenger `json:"passengers"`
	}

	var m Message

	var b testBus

	bus := testBus{
		[]testPassenger{
			{
				Name: "John",
				Age:  22,
			},
			{
				Name: "Steve",
				Age:  37,
			},
		},
	}

	data, _ := json.Marshal(bus)

	msg, _ := json.Marshal(NewMessage("peer", Transaction, data))
	_ = json.Unmarshal(msg, &m)
	_ = json.Unmarshal(m.Payload, &b)

	assert.Equal(t, bus, b)
}
