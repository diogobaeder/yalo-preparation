package clients

import (
	"errors"
	"github.com/nats-io/nats.go"
	"os"
)

// SimplifiedJetStream is a simplified version of nats.JetStreamContext containing only the methods we need.
// This will help with mocking the interface in tests.
type SimplifiedJetStream interface {
}

type Client struct {
	js SimplifiedJetStream
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
