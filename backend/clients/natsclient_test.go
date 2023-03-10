package clients

import (
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestStartsWithJetStream(t *testing.T) {
	require := require.New(t)
	if os.Getenv("NATS_URL") == "" {
		t.Skip("NATS_URL not configured, skipping the test.")
	}

	instance, err := NewNatsClient()

	require.Nilf(err, "Could not instantiate client: %v", err)
	require.Implements((*SimplifiedJetStream)(nil), instance.js)
	require.Implements((*nats.JetStreamContext)(nil), instance.js)
}
