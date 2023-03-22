package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"time"
	"yalo/diogo/demo/backend/internal/repositories"
)

type UserInfo struct {
	User string `uri:"user" binding:"required"`
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, "pong")
	})
	r.GET("/messages/latest-for/:user", func(c *gin.Context) {
		var info UserInfo
		if err := c.ShouldBindUri(&info); err != nil {
			c.JSON(400, gin.H{"msg": err})
			return
		}
		repo, err := repositories.NewMessagesRepository()
		if err != nil {
			c.JSON(500, gin.H{"msg": err})
			return
		}
		since := time.Now().UTC().Add(time.Duration(-24 * time.Hour))
		messages, err := repo.LatestForUser(info.User, since)
		if err != nil {
			c.JSON(500, gin.H{"msg": err})
			return
		}
		c.JSON(200, messages)
	})
	return r
}

func main() {
	port := os.Getenv("ADMIN_API_PORT")
	if port == "" {
		port = "8080"
	}
	r := setupRouter()
	_ = r.Run(fmt.Sprintf(":%v", port))
}
