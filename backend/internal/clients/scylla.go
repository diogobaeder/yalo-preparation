package clients

import (
	"errors"
	"github.com/gocql/gocql"
	"os"
	"strings"
)

type ScyllaClient struct {
	session *gocql.Session
}

func NewScyllaClient() (*ScyllaClient, error) {
	urls := os.Getenv("SCYLLA_ADDRS")
	if urls == "" {
		return &ScyllaClient{}, errors.New("SCYLLA_ADDRS environment variable not defined")
	}
	hosts := strings.Split(urls, ",")
	cluster := gocql.NewCluster(hosts...)
	session, err := cluster.CreateSession()
	if err != nil {
		return &ScyllaClient{}, errors.New("couldn't create session")
	}
	return &ScyllaClient{session}, nil
}
