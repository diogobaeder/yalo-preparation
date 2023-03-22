package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"yalo/diogo/demo/backend/internal/repositories"
)

type APISuite struct {
	suite.Suite
	router   *gin.Engine
	repo     *repositories.MessagesRepository
	recorder *httptest.ResponseRecorder
}

func (s *APISuite) SetupTest() {
	s.router = setupRouter()
	repo, err := repositories.NewMessagesRepository()
	s.Require().Nil(err)
	s.repo = repo
	s.recorder = httptest.NewRecorder()
	s.Require().Nil(repo.Truncate())
}

func (s *APISuite) TestPingRoute() {
	req, _ := http.NewRequest("GET", "/ping", nil)
	s.router.ServeHTTP(s.recorder, req)

	s.Require().Equal(200, s.recorder.Code)
	s.Require().Equal(`"pong"`, s.recorder.Body.String())
}

func (s *APISuite) TestLatestMessagesForUser() {
	message1 := repositories.NewMessage("johndoe", "foo", "request")
	time.Sleep(5 * time.Millisecond)
	message2 := repositories.NewMessage("johndoe", "bar", "reply")
	time.Sleep(5 * time.Millisecond)

	s.Require().Nil(s.repo.Insert(message1))
	s.Require().Nil(s.repo.Insert(message2))

	req, _ := http.NewRequest("GET", "/messages/latest-for/johndoe", nil)
	s.router.ServeHTTP(s.recorder, req)

	s.Require().Equal(200, s.recorder.Code)
	var messages []repositories.Message
	s.Require().Nil(json.Unmarshal(s.recorder.Body.Bytes(), &messages))
	s.Require().Equal(2, len(messages))
	s.Require().Equal([]repositories.Message{*message2, *message1}, messages)

}

func TestAPISuite(t *testing.T) {
	suite.Run(t, new(APISuite))
}
