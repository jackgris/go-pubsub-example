package main

import (
	"fmt"
	"sync"
)

// PubSub Interface: Start by defining an interface that represents the Pub/Sub functionality.
// This interface should include methods for publishing messages to topics and subscribing to topics.
type PubSub interface {
	Publish(topic string, message interface{})
	Subscribe(topic string) <-chan interface{}
}

// PubSubImpl implement the PubSub interface. This struct will manage the topics, subscribers, and message distribution.
type PubSubImpl struct {
	waitGroup        sync.WaitGroup
	topics           map[string][]chan interface{}
	subscriptionLock sync.Mutex
}

// NewPubSub create a struct that implements the PubSub interface.
func NewPubSub() *PubSubImpl {
	return &PubSubImpl{
		topics: make(map[string][]chan interface{}),
	}
}

// Publish in the PubSubImpl struct, implement the Publish method to
// publish messages to a specific topic. This method will iterate over the subscribers of the topic
// and send the message through the corresponding channels.
func (ps *PubSubImpl) Publish(topic string, message interface{}) {
	ps.subscriptionLock.Lock()
	defer ps.subscriptionLock.Unlock()

	subscribers := ps.topics[topic]
	for _, subscriber := range subscribers {
		ps.waitGroup.Add(1)
		go func(subscriber chan interface{}) {
			msg := fmt.Sprintf("%s %v", topic, message)
			subscriber <- msg
			ps.waitGroup.Done()
		}(subscriber)
	}
}

// Subscribe in the PubSubImpl implement the Subscribe method to allow subscribers to subscribe to a topic.
// It creates a new channel for the subscriber and adds it to the list of subscribers for the specified topic.
func (ps *PubSubImpl) Subscribe(topic string) <-chan interface{} {
	ps.subscriptionLock.Lock()
	defer ps.subscriptionLock.Unlock()

	subscriber := make(chan interface{})
	ps.topics[topic] = append(ps.topics[topic], subscriber)

	return subscriber
}

// Wait will wait until all messages are published
func (ps *PubSubImpl) Wait() {
	ps.waitGroup.Wait()
}
