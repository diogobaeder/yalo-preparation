package clients

import (
	"errors"
	"github.com/nats-io/nats.go"
	"os"
)

type Client struct {
	js nats.JetStreamContext
}

func NewNatsClient() (Client, error) {
	url := os.Getenv("NATS_URL")
	if url == "" {
		return Client{}, errors.New("NATS_URL not defined")
	}
	conn, err := nats.Connect(url)
	if err != nil {
		return Client{}, err
	}

	js, err := conn.JetStream()
	if err != nil {
		return Client{}, err
	}

	return Client{
		js,
	}, nil
}
