package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/nats-io/nats.go"
)

func main() {

	// Connect to a server
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		fmt.Println(err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	var shutdown sync.WaitGroup
	shutdown.Add(1)

	// Create JetStream Context
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		fmt.Println("Error jetstream subscriber: ", err)
	}

	// Create a Stream
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     "ORDERS",
		Subjects: []string{"ORDERS.*"},
	})
	if err != nil {
		fmt.Println("Error JetStream add stream: ", err)
	}

	ch := make(chan *nats.Msg, 64)

	// Simple Async Ephemeral Consumer
	_, err = js.Subscribe("ORDERS.*", func(m *nats.Msg) {
		ch <- m
	})
	if err != nil {
		fmt.Println("Error subscribe: ", err)
	}

	go func() {
		for {
			select {
			case msg := <-ch:
				fmt.Printf("Received a message: %s\n", string(msg.Data))
			case <-stop:
				shutdown.Done()
				break
			}
		}
	}()

	shutdown.Wait()
}
