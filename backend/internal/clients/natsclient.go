package clients

import (
	"errors"
	"github.com/nats-io/nats.go"
	"os"
	"time"
)

// SimplifiedJetStream is a simplified version of nats.JetStreamContext containing only the methods we need.
// This helps with mocking the interface in tests.
type SimplifiedJetStream interface {
	AddStream(cfg *nats.StreamConfig, opts ...nats.JSOpt) (*nats.StreamInfo, error)
	Subscribe(subj string, cb nats.MsgHandler, opts ...nats.SubOpt) (*nats.Subscription, error)
	QueueSubscribe(subj, queue string, cb nats.MsgHandler, opts ...nats.SubOpt) (*nats.Subscription, error)
	PublishAsync(subj string, data []byte, opts ...nats.PubOpt) (nats.PubAckFuture, error)
	PublishAsyncComplete() <-chan struct{}
}

// YaloNatsClient is a simplified facade to the JetStream client provided by NATS.
type YaloNatsClient struct {
	js SimplifiedJetStream
}

// Prepare makes sure the client is prepared to send and receive messages.
func (c *YaloNatsClient) Prepare() error {
	duration, _ := time.ParseDuration("24h")
	_, err := c.js.AddStream(&nats.StreamConfig{
		Name:     "yalo",
		Subjects: []string{"yalo.>"},
		MaxAge:   duration,
	})

	return err
}

// Subscribe asynchronously subscribes to a subject, without a queue group.
// Every subscriber defined like this will receive the messages for that subject.
func (c *YaloNatsClient) Subscribe(subj string, cb nats.MsgHandler) (*nats.Subscription, error) {
	subscription, err := c.js.Subscribe(subj, cb)
	return subscription, err
}

// QueueSubscribe asynchronously subscribe to a subject, but within a queue group.
// Only one subscriber for each queue group will receive the messages for that subject.
func (c *YaloNatsClient) QueueSubscribe(subj, queue string, cb nats.MsgHandler) (*nats.Subscription, error) {
	subscription, err := c.js.QueueSubscribe(subj, queue, cb)
	return subscription, err
}

// Publish asynchronously publishes a message to a subject.
// Make sure to call DonePublishing to make sure all messages were published successfully.
func (c *YaloNatsClient) Publish(subj string, data []byte) (nats.PubAckFuture, error) {
	paf, err := c.js.PublishAsync(subj, data)
	return paf, err
}

// DonePublishing returns a channel that can be used to check whether all messages have been published.
func (c *YaloNatsClient) DonePublishing() <-chan struct{} {
	return c.js.PublishAsyncComplete()
}

// NewNatsClient creates a new instance of YaloNatsClient, with sane defaults.
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
