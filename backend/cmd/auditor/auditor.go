package main

import (
	"github.com/nats-io/nats.go"
	"log"
	"yalo/diogo/demo/backend/internal/clients"
	"yalo/diogo/demo/backend/internal/repositories"
)

func main() {
	log.Println("Starting auditor...")
	client, err := clients.NewNatsClient()
	matcher := clients.NewSubjectMatcher()

	if err != nil {
		log.Panicf("Couldn't instantiate the client: %v", err)
	}

	if err = client.Prepare(); err != nil {
		log.Panicf("Couldn't prepare the client: %v", err)
	}

	repo, err := repositories.NewMessagesRepository()

	if err != nil {
		log.Panicf("Couldn't instantiate the repo: %v", err)
	}

	log.Println("Subscribing to subject within queue group...")
	_, err = client.QueueSubscribe("yalo.>", "auditors", func(msg *nats.Msg) {
		info := matcher.ExtractInfo(msg)
		log.Printf(`Got message from user %v: "%v"`, info.User, info.Message)
		message := repositories.NewMessage(info.User, info.Message, info.Direction)
		err := repo.Insert(message)
		if err != nil {
			log.Panicf("error inserting message in the database: %v", err)
		}
	})

	if err != nil {
		log.Panicf("Couldn't subscribe the client: %v", err)
	}
}
