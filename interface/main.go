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
	Wait()
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

var pubsub PubSub

// In this example, a message is published to "topic1", and the subscriber receives the message through the channel.
// You can have multiple subscribers to the same topic, and each subscriber will receive the published message independently.
// By leveraging channels and goroutines, you can achieve concurrent and asynchronous message distribution,
// enabling the Pub/Sub pattern in your Go application.
func main() {
	// Use the Pub/Sub implementation in your application by creating a new instance of PubSubImpl,
	// publishing messages, and subscribing to topics.
	pubsub = NewPubSub()

	// Subscribe to different topics
	subscriber1 := pubsub.Subscribe("topic1")
	subscriber2 := pubsub.Subscribe("topic2")
	subscriber3 := pubsub.Subscribe("topic3")
	subscriber4 := pubsub.Subscribe("topic3")
	subscriber5 := pubsub.Subscribe("topic3")

	// Publish a message to the topics
	pubsub.Publish("topic1", "Hello, subscribers number one!")
	pubsub.Publish("topic1", "Bye, subscribers number one!")
	pubsub.Publish("topic2", "Hello, subscribers number two!")
	pubsub.Publish("topic2", "How are you? subscribers number two!")
	pubsub.Publish("topic2", "Bye, subscribers number two!")
	pubsub.Publish("topic3", "Hello, subscribers number three!")
	pubsub.Publish("topic3", "How are you? subscribers number three!")
	pubsub.Publish("topic3", "Bye, subscribers number three!")

	// Receive messages from different topics
	go func() {
		for {
			select {
			case message := <-subscriber1:
				fmt.Println("subcriber 1", message)
			case message := <-subscriber2:
				fmt.Println("subcriber 2", message)
			case message := <-subscriber3:
				fmt.Println("subcriber 3", message)
			case message := <-subscriber4:
				fmt.Println("subcriber 4", message)
			case message := <-subscriber5:
				fmt.Println("subcriber 5", message)
			}

		}
	}()

	// Wait for all messages to be received by subscribers
	pubsub.Wait()
}
