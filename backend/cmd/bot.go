package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"regexp"
	"time"
	"yalo/diogo/demo/backend/internal/clients"
)

func main() {
	log.Println("Starting bot...")
	client, err := clients.NewNatsClient()
	botSubjectPattern := regexp.MustCompile(`^yalo\.bot\.(?P<user>[^.]+)$`)

	if err != nil {
		log.Panicf("Couldn't instantiate the client: %v", err)
	}

	if err = client.Prepare(); err != nil {
		log.Panicf("Couldn't prepare the client: %v", err)
	}

	log.Println("Subscribing to subject within queue group...")
	_, err = client.QueueSubscribe("yalo.bot.>", "bots", func(msg *nats.Msg) {
		matches := botSubjectPattern.FindStringSubmatch(msg.Subject)
		index := botSubjectPattern.SubexpIndex("user")
		user := matches[index]
		log.Printf("Got message fro user %v: %v", user, string(msg.Data))
		userSubject := fmt.Sprintf("yalo.user.%v", user)
		if _, err := client.Publish(userSubject, []byte("Got your message!")); err != nil {
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
			log.Panicln("Unable to finish publishing messages")
		}
	}
}