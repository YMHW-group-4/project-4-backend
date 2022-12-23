package networking

import "encoding/json"

// Message represents a message within the Network, whose
// singular purpose is to be spread to all connected peers.
type Message struct {
	Peer    string `json:"peer"`
	Topic   Topic  `json:"topic"`
	Payload string `json:"payload"`
}

// NewMessage creates a JSON encoded Message.
func NewMessage(peer string, topic Topic, payload string) ([]byte, error) {
	return json.Marshal(Message{
		Peer:    peer,
		Topic:   topic,
		Payload: payload,
	})
}
