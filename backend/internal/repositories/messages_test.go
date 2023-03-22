package repositories

import (
	"github.com/scylladb/gocqlx/v2"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type MessagesRepositorySuite struct {
	suite.Suite
	repo *MessagesRepository
}

func (s *MessagesRepositorySuite) SetupSuite() {
	repo, err := NewMessagesRepository()
	s.Require().Nil(err)
	s.repo = repo
	// Make sure we start with a blank slate
	s.Require().Nil(s.repo.Truncate())
}

func (s *MessagesRepositorySuite) TestStartsWithSession() {
	s.Require().IsType(&gocqlx.Session{}, s.repo.session)
}

func (s *MessagesRepositorySuite) TestQueriesLatestMessagesForUser() {
	message := NewMessage("johndoe", "Something to say", "request")
	s.Require().Nil(s.repo.Insert(message))

	retrieved, err := s.repo.LatestForUser("johndoe", message.Time.Add(time.Hour*-1))

	s.Require().Nil(err)
	s.Require().Equal(1, len(retrieved))
	s.Require().Equal(Message{
		User:      "johndoe",
		Time:      message.Time,
		Message:   "Something to say",
		Direction: "request",
	}, *retrieved[0])
}

func TestMessagesRepositorySuite(t *testing.T) {
	suite.Run(t, new(MessagesRepositorySuite))
}

type MessagesSuite struct {
	suite.Suite
}

func (s *MessagesSuite) TestCreatesNewRequest() {
	now := time.UnixMilli(time.Now().UnixMilli()).UTC()

	message := NewMessage("johndoe", "Something useful", "request")

	s.Require().Equal("johndoe", message.User)
	s.Require().Equal("Something useful", message.Message)
	s.Require().Equal("request", message.Direction)
	s.Require().GreaterOrEqual(message.Time, now)
}

func TestMessageSuite(t *testing.T) {
	suite.Run(t, new(MessagesSuite))
}
