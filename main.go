package main

import (
	"fmt"
	"sync"
)

// Create a goroutine that will publish messages to the channel.
func publisher(wg *sync.WaitGroup, msgChan chan string) {
	for i := 0; i < 10; i++ {
		msgChan <- fmt.Sprintf("Message %d", i)
	}
	close(msgChan)
	wg.Done()
}

// Create one or more goroutines that will subscribe to the channel and receive the published messages.
func subscriber(id int, wg *sync.WaitGroup, msgChan chan string) {
	for message := range msgChan {
		fmt.Printf("Subscriber %d received message: %s\n", id, message)
	}
	wg.Done()
}

// In the main function, start the publisher and subscriber goroutines.
func main() {
	// Define a channel that will be used to communicate between publishers and subscribers.
	var msgChan = make(chan string)
	// We use a WaitGroup to coordinate goroutines
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go publisher(wg, msgChan)
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go subscriber(i, wg, msgChan)
	}
	// Wait for all the subscribers to finish receiving messages.
	wg.Wait()
}
