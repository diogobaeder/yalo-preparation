package clients

import (
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"os"
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

func TestStartsWithJetStream(t *testing.T) {
	ensure := require.New(t)
	if os.Getenv("NATS_URL") == "" {
		t.Skip("NATS_URL not configured, skipping the test.")
	}

	instance, err := NewNatsClient()

	ensure.Nilf(err, "Could not instantiate client: %v", err)
	ensure.Implements((*SimplifiedJetStream)(nil), instance.js)
	ensure.Implements((*nats.JetStreamContext)(nil), instance.js)
}

func TestPreparesClient(t *testing.T) {
	ensure := require.New(t)
	js := new(JetStreamContextMock)
	instance := &Client{js}
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
