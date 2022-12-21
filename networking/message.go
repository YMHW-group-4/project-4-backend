package networking

import "github.com/libp2p/go-libp2p/core/peer"

type Message struct {
	Peer    peer.ID
	Topic   Topic
	Payload string
}

func NewMessage(peer peer.ID, topic Topic, payload string) Message {
	return Message{
		Peer:    peer,
		Topic:   topic,
		Payload: payload,
	}
}
