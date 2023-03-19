package clients

import (
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
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

type NatsClientSuite struct {
	suite.Suite
	client *NatsClient
	js     *JetStreamContextMock
}

func (s *NatsClientSuite) SetupSuite() {
	js := new(JetStreamContextMock)
	client := &NatsClient{js}
	s.js = js
	s.client = client
}

func (s *NatsClientSuite) TestStartsWithJetStream() {
	instance, err := NewNatsClient()

	s.Require().Nilf(err, "Could not instantiate client: %v", err)
	s.Require().Implements((*SimplifiedJetStream)(nil), instance.js)
	s.Require().Implements((*nats.JetStreamContext)(nil), instance.js)
}

func (s *NatsClientSuite) TestPreparesClient() {
	duration, _ := time.ParseDuration("24h")
	config := &nats.StreamConfig{
		Name:     "yalo",
		Subjects: []string{"yalo.>"},
		MaxAge:   duration,
	}
	info := new(nats.StreamInfo)
	s.js.On("AddStream", config).Return(info, nil)

	_ = s.client.Prepare()

	s.Require().True(s.js.AssertCalled(s.T(), "AddStream", config))
}

func (s *NatsClientSuite) TestSubscribesToSubject() {
	subject := "yalo.something"
	callback := func(msg *nats.Msg) {}
	subscription := new(nats.Subscription)
	// Note: unfortunately I can only make this test pass if I use mock.Anything to match the callback function.
	// In the future hopefully I can use testify's more specific matchers.
	s.js.On("Subscribe", subject, mock.Anything).Return(subscription, nil)

	_, err := s.client.Subscribe(subject, callback)

	s.Require().Nil(err)
	s.Require().True(s.js.AssertCalled(s.T(), "Subscribe", subject, mock.Anything))
}

func (s *NatsClientSuite) TestSubscribesToSubjectInQueue() {
	subject := "yalo.something"
	queue := "some_queue"
	callback := func(msg *nats.Msg) {}
	subscription := new(nats.Subscription)
	// Note: unfortunately I can only make this test pass if I use mock.Anything to match the callback function.
	// In the future hopefully I can use testify's more specific matchers.
	s.js.On("QueueSubscribe", subject, queue, mock.Anything).Return(subscription, nil)

	_, err := s.client.QueueSubscribe(subject, queue, callback)

	s.Require().Nil(err)
	s.Require().True(s.js.AssertCalled(s.T(), "QueueSubscribe", subject, queue, mock.Anything))
}

func (s *NatsClientSuite) TestPublishesToSubject() {
	subject := "yalo.something"
	data := []byte("somewhere")
	paf := new(pubAckFuture)
	s.js.On("PublishAsync", subject, data).Return(paf, nil)

	_, err := s.client.Publish(subject, data)

	s.Require().Nil(err)
	s.Require().True(s.js.AssertCalled(s.T(), "PublishAsync", subject, data))
}

func (s *NatsClientSuite) TestChecksIsDonePublishing() {
	channel := make(<-chan struct{})
	s.js.On("PublishAsyncComplete").Return(channel)

	s.client.DonePublishing()

	s.Require().True(s.js.AssertCalled(s.T(), "PublishAsyncComplete"))
}

func TestNatsSuite(t *testing.T) {
	suite.Run(t, new(NatsClientSuite))
}

type SubjectMatcherSuite struct {
	suite.Suite
	matcher *SubjectMatcher
}

func (s *SubjectMatcherSuite) SetupSuite() {
	s.matcher = NewSubjectMatcher()
}

func (s *SubjectMatcherSuite) TestFindsUser() {
	s.Require().Equal(s.matcher.FindUser("yalo.bot.johndoe"), "johndoe")
}

func (s *SubjectMatcherSuite) TestExtractsUserInfoFromNatsMessage() {
	msg := &nats.Msg{
		Subject: "yalo.bot.johndoe",
		Data:    []byte("Just said something"),
	}

	info := s.matcher.ExtractInfo(msg)

	s.Require().Equal(info.User, "johndoe")
	s.Require().Equal(info.Message, "Just said something")
	s.Require().Equal(info.ReplyTo, "yalo.user.johndoe")
}

func TestSubjectMatcherSuite(t *testing.T) {
	suite.Run(t, new(SubjectMatcherSuite))
}
