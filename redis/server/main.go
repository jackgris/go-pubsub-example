package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	ctx := context.Background()

	// There is no error because go-redis automatically reconnects on error.
	pubsub := rdb.Subscribe(ctx, "send-user-data")
	// Close the subscription when we are done.
	defer pubsub.Close()

	ch := pubsub.Channel()

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello there 👋")
	})

	app.Post("/", func(c *fiber.Ctx) error {
		user := new(User)

		if err := c.BodyParser(user); err != nil {
			log.Println("Body Parse: ", err)
			return err
		}

		payload, err := json.Marshal(user)
		if err != nil {
			return err
		}

		if err := rdb.Publish(ctx, "send-user-data", payload).Err(); err != nil {
			return err
		}

		return c.SendStatus(200)
	})

	go func(ch <-chan *redis.Message) {
		for msg := range ch {
			fmt.Println(msg.Channel, msg.Payload)
		}
	}(ch)

	log.Println(app.Listen(":3000"))
}
