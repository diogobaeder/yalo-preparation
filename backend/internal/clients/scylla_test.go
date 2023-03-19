package clients

import (
	"github.com/scylladb/gocqlx/v2"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type ClientSuite struct {
	suite.Suite
	client *ScyllaClient
}

func (s *ClientSuite) SetupSuite() {
	client, err := NewScyllaClient()
	s.Nil(err)
	s.client = client
	// Make sure we start with a blank slate
	s.Require().Nil(s.client.Truncate())
}

func (s *ClientSuite) TestScyllaClientStartsWithSession() {
	s.Require().IsType(&gocqlx.Session{}, s.client.session)
}

func (s *ClientSuite) TestInsertsMessage() {
	now := time.UnixMilli(time.Now().UnixMilli()).UTC()
	past := now.Add(time.Hour * -1)
	message := &Message{
		"johndoe",
		now,
		"Something to say",
	}

	err := s.client.Insert(message)

	s.Require().Nilf(err, "Could not insert message: %v", err)

	retrieved, err := s.client.LatestForUser("johndoe", past)

	s.Require().Nilf(err, "Could not retrieve messages: %v", err)

	s.Require().Equal(1, len(retrieved))
	s.Require().Equal(Message{
		User:    "johndoe",
		Time:    now,
		Message: "Something to say",
	}, *retrieved[0])
}

type MessagesRepositorySuite struct {
	suite.Suite
	repo *MessagesRepository
}

func (s *MessagesRepositorySuite) SetupSuite() {
	repo, err := NewMessagesRepository()
	s.Require().Nil(err)
	s.repo = repo
	s.Require().Nil(s.repo.truncate())
}

func (s *MessagesRepositorySuite) TestMessagesRepositoryStartsWithScyllaClient() {
	s.Require().IsType(&ScyllaClient{}, s.repo.client)
}

func (s *MessagesRepositorySuite) TestMessagesRepositoryQueriesLatestMessagesForUser() {
	s.Require().IsType(&ScyllaClient{}, s.repo.client)
}

func TestClientSuite(t *testing.T) {
	suite.Run(t, new(ClientSuite))
}

func TestMessagesRepositorySuite(t *testing.T) {
	suite.Run(t, new(MessagesRepositorySuite))
}
