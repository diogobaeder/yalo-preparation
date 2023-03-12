package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"time"
	"yalo/diogo/demo/backend/internal/clients"
)

func main() {
	log.Println("Starting bot...")
	client, err := clients.NewNatsClient()

	if err != nil {
		log.Panicf("Couldn't instantiate the client: %v", err)
	}

	err = client.Prepare()

	if err != nil {
		log.Panicf("Couldn't prepare the client: %v", err)
	}

	log.Println("Subscribing to subject within queue group...")
	_, err = client.QueueSubscribe("yalo.bot.>", "bots", func(msg *nats.Msg) {
		log.Printf("Got message: %v", string(msg.Data))
		_, err := client.Publish("yalo.user.1234", []byte("Got your message!"))

		if err != nil {
			log.Panicf("Couldn't publish message: %v", err)
		}
	})

	if err != nil {
		log.Panicf("Couldn't subscribe the client: %v", err)
	}

	log.Println("Now waiting for messages.")
	for {
		select {
		case <-client.DonePublishing():
		case <-time.After(5 * time.Second):
			fmt.Println("Unable to finish publishing messages")
		}
	}
}
