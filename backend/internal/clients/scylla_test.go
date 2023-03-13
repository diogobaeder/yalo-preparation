package clients

import (
	"github.com/scylladb/gocqlx/v2"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestStartsWithSession(t *testing.T) {
	ensure := require.New(t)
	if os.Getenv("SCYLLA_ADDRS") == "" {
		t.Skip("SCYLLA_ADDRS not configured, skipping the test.")
	}

	instance, err := NewScyllaClient()

	ensure.Nilf(err, "Could not instantiate client: %v", err)
	ensure.IsType(&gocqlx.Session{}, instance.session)
}
