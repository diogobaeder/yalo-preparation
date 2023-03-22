package clients

import (
	"errors"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"os"
	"regexp"
	"time"
)

const STREAM_NAMESPACE string = "yalo"
const STREAM_DURATION string = "24h"
const STREAM_SUBJECTS string = "yalo.>"

// SimplifiedJetStream is a simplified version of nats.JetStreamContext containing only the methods we need.
// This helps with mocking the interface in tests.
type SimplifiedJetStream interface {
	AddStream(cfg *nats.StreamConfig, opts ...nats.JSOpt) (*nats.StreamInfo, error)
	Subscribe(subj string, cb nats.MsgHandler, opts ...nats.SubOpt) (*nats.Subscription, error)
	QueueSubscribe(subj, queue string, cb nats.MsgHandler, opts ...nats.SubOpt) (*nats.Subscription, error)
	PublishAsync(subj string, data []byte, opts ...nats.PubOpt) (nats.PubAckFuture, error)
	PublishAsyncComplete() <-chan struct{}
}

// NatsClient is a simplified facade to the JetStream client provided by NATS.
type NatsClient struct {
	js SimplifiedJetStream
}

// Prepare makes sure the client is prepared to send and receive messages.
func (c *NatsClient) Prepare() error {
	duration, _ := time.ParseDuration(STREAM_DURATION)
	_, err := c.js.AddStream(&nats.StreamConfig{
		Name:     STREAM_NAMESPACE,
		Subjects: []string{STREAM_SUBJECTS},
		MaxAge:   duration,
	})

	return err
}

// Subscribe asynchronously subscribes to a subject, without a queue group.
// Every subscriber defined like this will receive the messages for that subject.
func (c *NatsClient) Subscribe(subj string, cb nats.MsgHandler) (*nats.Subscription, error) {
	subscription, err := c.js.Subscribe(subj, cb)
	return subscription, err
}

// QueueSubscribe asynchronously subscribe to a subject, but within a queue group.
// Only one subscriber for each queue group will receive the messages for that subject.
func (c *NatsClient) QueueSubscribe(subj, queue string, cb nats.MsgHandler) (*nats.Subscription, error) {
	subscription, err := c.js.QueueSubscribe(subj, queue, cb)
	return subscription, err
}

// Publish asynchronously publishes a message to a subject.
// Make sure to call DonePublishing to make sure all messages were published successfully.
func (c *NatsClient) Publish(subj string, data []byte) (nats.PubAckFuture, error) {
	paf, err := c.js.PublishAsync(subj, data)
	return paf, err
}

// DonePublishing returns a channel that can be used to check whether all messages have been published.
func (c *NatsClient) DonePublishing() <-chan struct{} {
	return c.js.PublishAsyncComplete()
}

// NewNatsClient creates a new instance of NatsClient, with sane defaults.
func NewNatsClient() (*NatsClient, error) {
	url := os.Getenv("NATS_URL")
	if url == "" {
		return nil, errors.New("NATS_URL environment variable not defined")
	}
	log.Printf("NATS_URL: %v\n", url)
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	js, err := conn.JetStream()
	if err != nil {
		return nil, err
	}

	return &NatsClient{
		js,
	}, nil
}

type UserInfo struct {
	User      string
	Message   string
	ReplyTo   string
	Direction string
}

type SubjectMatcher struct {
	subjectPattern *regexp.Regexp
}

func (s *SubjectMatcher) FindUser(subject string) string {
	matches := s.subjectPattern.FindStringSubmatch(subject)
	index := s.subjectPattern.SubexpIndex("user")
	return matches[index]
}

func (s *SubjectMatcher) FindDirection(subject string) string {
	matches := s.subjectPattern.FindStringSubmatch(subject)
	index := s.subjectPattern.SubexpIndex("direction")
	return matches[index]
}

func (s *SubjectMatcher) ExtractInfo(msg *nats.Msg) *UserInfo {
	user := s.FindUser(msg.Subject)
	direction := s.FindDirection(msg.Subject)
	return &UserInfo{
		User:      user,
		Message:   string(msg.Data),
		ReplyTo:   fmt.Sprintf(`%v.reply.%v`, STREAM_NAMESPACE, user),
		Direction: direction,
	}
}

func NewSubjectMatcher() *SubjectMatcher {
	return &SubjectMatcher{
		subjectPattern: regexp.MustCompile(fmt.Sprintf(`^%v\.(?P<direction>request|reply)\.(?P<user>[^.]+)$`, STREAM_NAMESPACE)),
	}
}
