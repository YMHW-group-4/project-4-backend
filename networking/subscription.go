package networking

import (
	"context"
	"encoding/json"
	"sync"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"
)

// source: https://medium.com/rahasak/libp2p-pubsub-peer-discovery-with-kademlia-dht-c8b131550ac7

// Topic the Topic to whom a Subscription can be made.
type Topic string

const (
	Transaction Topic = "transaction"
	Block       Topic = "block"
	Blockchain  Topic = "blockchain"
	Consensus   Topic = "consensus"
)

// Subscription represents a Subscription within the Network.
type Subscription struct {
	Messages chan Message
	topic    Topic
	self     peer.ID
	ctx      context.Context
	top      *pubsub.Topic
	sub      *pubsub.Subscription
	wg       sync.WaitGroup
	close    chan struct{}
}

// NewSubscription creates a new Subscription on given topic.
func NewSubscription(ctx context.Context, ps *pubsub.PubSub, host peer.ID, topic Topic) (*Subscription, error) {
	top, err := ps.Join(string(topic))
	if err != nil {
		return nil, err
	}

	sub, err := top.Subscribe()
	if err != nil {
		return nil, err
	}

	s := &Subscription{
		Messages: make(chan Message, 0),
		topic:    topic,
		self:     host,
		ctx:      ctx,
		top:      top,
		sub:      sub,
		wg:       sync.WaitGroup{},
		close:    make(chan struct{}, 0),
	}

	s.listen()

	return s, nil
}

// Publish publishes a Message to this Subscription.
func (subscription *Subscription) Publish(message []byte) error {
	return subscription.top.Publish(subscription.ctx, message)
}

// Close closes a Subscription.
func (subscription *Subscription) Close() error {
	close(subscription.close)
	subscription.sub.Cancel()

	if err := subscription.top.Close(); err != nil {
		return err
	}

	subscription.wg.Wait()

	return nil
}

// listen listens to incoming Messages.
func (subscription *Subscription) listen() {
	subscription.wg.Add(2) //nolint

	ch := make(chan *pubsub.Message, 0)

	go func() {
		defer subscription.wg.Done()

		for {
			msg, err := subscription.sub.Next(subscription.ctx)
			if err != nil {
				return
			}

			if msg.ReceivedFrom != subscription.self {
				ch <- msg
			}
		}
	}()

	go func() {
		defer subscription.wg.Done()

		for {
			select {
			case <-subscription.close:
				return
			case msg := <-ch:
				var message Message

				if err := json.Unmarshal(msg.Data, &message); err == nil {
					subscription.Messages <- message
				}
			}
		}
	}()
}
