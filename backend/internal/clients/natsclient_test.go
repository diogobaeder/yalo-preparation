package clients

import (
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type JetStreamContextMock struct {
	mock.Mock
}

func (j *JetStreamContextMock) AddStream(cfg *nats.StreamConfig, opts ...nats.JSOpt) (*nats.StreamInfo, error) {
	_ = opts
	args := j.Called(cfg)
	return args.Get(0).(*nats.StreamInfo), nil
}

func (j *JetStreamContextMock) Subscribe(subj string, cb nats.MsgHandler, opts ...nats.SubOpt) (*nats.Subscription, error) {
	_ = opts
	args := j.Called(subj, cb)
	return args.Get(0).(*nats.Subscription), nil
}

func (j *JetStreamContextMock) QueueSubscribe(subj, queue string, cb nats.MsgHandler, opts ...nats.SubOpt) (*nats.Subscription, error) {
	_ = opts
	args := j.Called(subj, queue, cb)
	return args.Get(0).(*nats.Subscription), nil
}

func (j *JetStreamContextMock) PublishAsync(subj string, data []byte, opts ...nats.PubOpt) (nats.PubAckFuture, error) {
	_ = opts
	args := j.Called(subj, data)
	return args.Get(0).(nats.PubAckFuture), nil
}

func (j *JetStreamContextMock) PublishAsyncComplete() <-chan struct{} {
	args := j.Called()
	return args.Get(0).(<-chan struct{})
}

type pubAckFuture struct{}

func (p *pubAckFuture) Ok() <-chan *nats.PubAck {
	channel := make(chan *nats.PubAck)
	return channel
}

func (p *pubAckFuture) Err() <-chan error {
	channel := make(chan error)
	return channel
}

func (p *pubAckFuture) Msg() *nats.Msg {
	return new(nats.Msg)
}

func TestStartsWithJetStream(t *testing.T) {
	ensure := require.New(t)

	instance, err := NewNatsClient()

	ensure.Nilf(err, "Could not instantiate client: %v", err)
	ensure.Implements((*SimplifiedJetStream)(nil), instance.js)
	ensure.Implements((*nats.JetStreamContext)(nil), instance.js)
}

func TestPreparesClient(t *testing.T) {
	ensure := require.New(t)
	js := new(JetStreamContextMock)
	instance := &NatsClient{js}
	duration, _ := time.ParseDuration("24h")
	config := &nats.StreamConfig{
		Name:     "yalo",
		Subjects: []string{"yalo.>"},
		MaxAge:   duration,
	}
	info := new(nats.StreamInfo)
	js.On("AddStream", config).Return(info, nil)

	_ = instance.Prepare()

	ensure.True(js.AssertCalled(t, "AddStream", config))
}

func TestSubscribesToSubject(t *testing.T) {
	ensure := require.New(t)
	js := new(JetStreamContextMock)
	instance := &NatsClient{js}
	subject := "yalo.something"
	callback := func(msg *nats.Msg) {}
	subscription := new(nats.Subscription)
	// Note: unfortunately I can only make this test pass if I use mock.Anything to match the callback function.
	// In the future hopefully I can use testify's more specific matchers.
	js.On("Subscribe", subject, mock.Anything).Return(subscription, nil)

	_, err := instance.Subscribe(subject, callback)

	ensure.Nil(err)
	ensure.True(js.AssertCalled(t, "Subscribe", subject, mock.Anything))
}

func TestSubscribesToSubjectInQueue(t *testing.T) {
	ensure := require.New(t)
	js := new(JetStreamContextMock)
	instance := &NatsClient{js}
	subject := "yalo.something"
	queue := "some_queue"
	callback := func(msg *nats.Msg) {}
	subscription := new(nats.Subscription)
	// Note: unfortunately I can only make this test pass if I use mock.Anything to match the callback function.
	// In the future hopefully I can use testify's more specific matchers.
	js.On("QueueSubscribe", subject, queue, mock.Anything).Return(subscription, nil)

	_, err := instance.QueueSubscribe(subject, queue, callback)

	ensure.Nil(err)
	ensure.True(js.AssertCalled(t, "QueueSubscribe", subject, queue, mock.Anything))
}

func TestPublishesToSubject(t *testing.T) {
	ensure := require.New(t)
	js := new(JetStreamContextMock)
	instance := &NatsClient{js}
	subject := "yalo.something"
	data := []byte("somewhere")
	paf := new(pubAckFuture)
	js.On("PublishAsync", subject, data).Return(paf, nil)

	_, err := instance.Publish(subject, data)

	ensure.Nil(err)
	ensure.True(js.AssertCalled(t, "PublishAsync", subject, data))
}

func TestChecksIsDonePublishing(t *testing.T) {
	ensure := require.New(t)
	js := new(JetStreamContextMock)
	instance := &NatsClient{js}
	channel := make(<-chan struct{})
	js.On("PublishAsyncComplete").Return(channel)

	instance.DonePublishing()

	ensure.True(js.AssertCalled(t, "PublishAsyncComplete"))
}
