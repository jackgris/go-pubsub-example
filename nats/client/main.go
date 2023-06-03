package main

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	// Connect to NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		fmt.Println("Error while connect: ", err)
	}

	// Create JetStream Context
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		fmt.Println("Error JetStream Context: ", err)
	}

	// Simple Async Stream Publisher
	for i := 0; i < 500; i++ {
		hello := fmt.Sprintf("Hello number %d", i)
		_, err = js.PublishAsync("ORDERS.scratch", []byte(hello))
		if err != nil {
			fmt.Println("Error publishing: ", err)
		}
	}

	select {
	case <-js.PublishAsyncComplete():
	case <-time.After(5 * time.Second):
		fmt.Println("Did not resolve in time")
	}
}
