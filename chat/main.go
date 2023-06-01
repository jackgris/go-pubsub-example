package main

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

func main() {
	// Create our server and configure the middleware that will help us to establish the WebSocket connection.
	app := fiber.New()
	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// Initiate the Pub/Sub service created by us.
	var pubsub PubSub = NewPubSub()
	// In this case we only have one topic, one chat room.
	subscribe := pubsub.Subscribe("chat")
	// With this will save all the websocket connections that are active.
	connections := make(map[string]*websocket.Conn)

	app.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {
		// c.Locals is added to the *websocket.Conn
		log.Println(c.Locals("allowed"))  // true
		log.Println(c.Params("id"))       // 123
		log.Println(c.Query("v"))         // 1.0
		log.Println(c.Cookies("session")) // ""

		// Create an Id to identify the connection.
		id := uuid.New()
		// Save the connection in our pool. And remove the connection when the websocket is close.
		connections[id.String()] = c
		defer delete(connections, id.String())

		var (
			mt  int
			msg []byte
			err error
		)
		for {
			// We read everything we receive from the front-end.
			if mt, msg, err = c.ReadMessage(); err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s - mt: %d", msg, mt)

			// We get the data from the message and publish that in our Pub/Sub service.
			m := Msg{}
			err := json.Unmarshal(msg, &m)
			if err != nil {
				log.Println(err)
				continue
			}
			m.Id = id.String() // add the Id to identify the owner
			pubsub.Publish("chat", m)
		}

	}))

	// Every time that receive a message, we will send the message to all the active connections.
	// For that we need the channel for communication and all the connections.
	go func(subs <-chan Msg, conn map[string]*websocket.Conn) {
		for {
			message := <-subs
			m, err := json.Marshal(message)
			if err != nil {
				log.Println("Error while marshal our message: ", err)
				continue
			}

			for _, c := range conn {
				if err = c.WriteMessage(1, m); err != nil {
					log.Printf("Error while write to ID: %s with err: %s\n", message.Id, err)
				}
			}
		}
	}(subscribe, connections)

	log.Fatal(app.Listen(":8080"))
	// Access the websocket server: ws://localhost:3000/ws/123?v=1.0
	// https://www.websocket.org/echo.html
}
