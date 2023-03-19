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
	s.Require().Nil(s.repo.truncate())
}

func (s *MessagesRepositorySuite) TestStartsWithSession() {
	s.Require().IsType(&gocqlx.Session{}, s.repo.session)
}

func (s *MessagesRepositorySuite) TestQueriesLatestMessagesForUser() {
	now := time.UnixMilli(time.Now().UnixMilli()).UTC()
	past := now.Add(time.Hour * -1)
	message := &Message{
		"johndoe",
		now,
		"Something to say",
	}

	s.Require().Nil(s.repo.Insert(message))
	retrieved, err := s.repo.LatestForUser("johndoe", past)

	s.Require().Nil(err)

	s.Require().Equal(1, len(retrieved))
	s.Require().Equal(Message{
		User:    "johndoe",
		Time:    now,
		Message: "Something to say",
	}, *retrieved[0])
}

func TestMessagesRepositorySuite(t *testing.T) {
	suite.Run(t, new(MessagesRepositorySuite))
}
