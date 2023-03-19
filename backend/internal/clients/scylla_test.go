package clients

import (
	"github.com/scylladb/gocqlx/v2"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
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

func TestInsertsMessage(t *testing.T) {
	ensure := require.New(t)
	if os.Getenv("SCYLLA_ADDRS") == "" {
		t.Skip("SCYLLA_ADDRS not configured, skipping the test.")
	}
	instance, _ := NewScyllaClient()
	now := time.UnixMilli(time.Now().UnixMilli()).UTC()
	past := now.Add(time.Hour * -1)
	message := &Message{
		"johndoe",
		now,
		"Something to say",
	}
	ensure.Nil(instance.Truncate())

	err := instance.Insert(message)

	ensure.Nilf(err, "Could not insert message: %v", err)

	retrieved, err := instance.LatestForUser("johndoe", past)

	ensure.Nilf(err, "Could not retrieve messages: %v", err)

	ensure.Equal(1, len(retrieved))
	ensure.Equal(Message{
		User:    "johndoe",
		Time:    now,
		Message: "Something to say",
	}, *retrieved[0])
}
