package clients

import (
	"github.com/gocql/gocql"
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
	ensure.IsType((*gocql.Session)(nil), instance.session)
}
