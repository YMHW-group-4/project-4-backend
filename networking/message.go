package networking

import "github.com/libp2p/go-libp2p/core/peer"

type Message struct {
	Peer    peer.ID
	Payload string `json:"payload"`
}

func NewMessage(peer peer.ID, payload string) Message {
	return Message{
		Peer:    peer,
		Payload: payload,
	}
}
