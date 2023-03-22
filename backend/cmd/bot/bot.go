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
	matcher := clients.NewSubjectMatcher()

	if err != nil {
		log.Panicf("Couldn't instantiate the client: %v", err)
	}

	if err = client.Prepare(); err != nil {
		log.Panicf("Couldn't prepare the client: %v", err)
	}

	log.Println("Subscribing to subject within queue group...")
	_, err = client.QueueSubscribe("yalo.request.>", "bots", func(msg *nats.Msg) {
		info := matcher.ExtractInfo(msg)
		log.Printf(`Got message from user %v: "%v"`, info.User, info.Message)
		botMessage := fmt.Sprintf(`Got your message, %v! This is what you said: "%v"`, info.User, info.Message)
		if _, err := client.Publish(info.ReplyTo, []byte(botMessage)); err != nil {
			log.Panicf("Couldn't publish message: %v", err)
		}
	})

	if err != nil {
		log.Panicf("Couldn't subscribe the client: %v", err)
	}

	log.Println("Now waiting for messages.")
	for {
		select {
		default:
			time.Sleep(1 * time.Millisecond)
		}
	}
}
