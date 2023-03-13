package clients

import (
	"errors"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
	"os"
	"strings"
)

type ScyllaClient struct {
	session *gocqlx.Session
}

func NewScyllaClient() (*ScyllaClient, error) {
	urls := os.Getenv("SCYLLA_ADDRS")
	if urls == "" {
		return &ScyllaClient{}, errors.New("SCYLLA_ADDRS environment variable not defined")
	}
	hosts := strings.Split(urls, ",")
	cluster := gocql.NewCluster(hosts...)
	session, err := gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		return &ScyllaClient{}, errors.New("couldn't create session")
	}
	return &ScyllaClient{&session}, nil
}
