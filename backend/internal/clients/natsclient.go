package clients

import (
	"errors"
	"github.com/nats-io/nats.go"
	"os"
	"time"
)

// SimplifiedJetStream is a simplified version of nats.JetStreamContext containing only the methods we need.
// This will help with mocking the interface in tests.
type SimplifiedJetStream interface {
	AddStream(cfg *nats.StreamConfig, opts ...nats.JSOpt) (*nats.StreamInfo, error)
}

type YaloNatsClient struct {
	js SimplifiedJetStream
}

func (c *YaloNatsClient) Prepare() error {
	duration, _ := time.ParseDuration("24h")
	_, err := c.js.AddStream(&nats.StreamConfig{
		Name:     "yalo",
		Subjects: []string{"yalo.>"},
		MaxAge:   duration,
	})

	return err
}

func NewNatsClient() (YaloNatsClient, error) {
	url := os.Getenv("NATS_URL")
	if url == "" {
		return YaloNatsClient{}, errors.New("NATS_URL not defined")
	}
	conn, err := nats.Connect(url)
	if err != nil {
		return YaloNatsClient{}, err
	}

	js, err := conn.JetStream()
	if err != nil {
		return YaloNatsClient{}, err
	}

	return YaloNatsClient{
		js,
	}, nil
}
