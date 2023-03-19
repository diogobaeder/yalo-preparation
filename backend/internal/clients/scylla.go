package clients

import (
	"errors"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/qb"
	"github.com/scylladb/gocqlx/v2/table"
	"os"
	"strings"
	"time"
)

var messageMetadata = table.Metadata{
	Name:    "messages",
	Columns: []string{"user", "time", "message"},
	PartKey: []string{"user"},
	SortKey: []string{"time"},
}

var messageTable = table.New(messageMetadata)

type Message struct {
	User    string
	Time    time.Time
	Message string
}

type ScyllaClient struct {
	session *gocqlx.Session
}

func (c *ScyllaClient) Insert(message *Message) error {
	query := c.session.Query(messageTable.Insert()).BindStruct(message)
	return query.ExecRelease()
}

func (c *ScyllaClient) Truncate() error {
	return c.session.ExecStmt(fmt.Sprintf("TRUNCATE %v", messageTable.Name()))
}

func (c *ScyllaClient) LatestForUser(user string, since time.Time) ([]*Message, error) {
	var messages []*Message
	//query := c.session.Query(messageTable.Select()).BindMap(qb.M{
	//	"user": user,
	//	"time": qb.GtOrEqLit("time", since.String()),
	//})
	query := c.session.Query(qb.Select(messageTable.Name()).Where(
		qb.EqLit("user", fmt.Sprintf("'%v'", user)),
		qb.GtOrEqLit("time", since.Format("'2006-01-02 15:04:05.999'")),
	).ToCql())
	//log.Fatalln(query.String())
	if err := query.SelectRelease(&messages); err != nil {
		return nil, err
	}
	return messages, nil
}

func NewScyllaClient() (*ScyllaClient, error) {
	urls := os.Getenv("SCYLLA_ADDRS")
	keyspace := os.Getenv("SCYLLA_KEYSPACE")
	if urls == "" || keyspace == "" {
		return &ScyllaClient{}, errors.New("both SCYLLA_ADDRS and SCYLLA_KEYSPACE vars should be defined")
	}
	hosts := strings.Split(urls, ",")
	cluster := gocql.NewCluster(hosts...)
	cluster.Keyspace = keyspace
	session, err := gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		return &ScyllaClient{}, err
	}
	return &ScyllaClient{&session}, nil
}
