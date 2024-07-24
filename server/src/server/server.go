package server

import (
	"net/http"
	"os"

	"github.com/ataboo/rtc-game-buzzer/src/room"
	"github.com/gin-gonic/gin"
)

type JoinInput struct {
	RoomCode string
	Name     string
}

type HostInput struct {
	Name string
}

func Start() {
	hostAuthKey := os.Getenv("HOST_AUTH_KEY")
	addr := os.Getenv("HOSTNAME")

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	r.POST("/host", func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader != hostAuthKey {
			c.JSON(http.StatusForbidden, gin.H{"message": "Not authorized"})
			return
		}

		input := HostInput{}
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
			return
		}

		status := room.CreateRoom(c, input.Name)
		if status != http.StatusOK {
			c.JSON(status, gin.H{"message": "Failed to host"})
		}
	})

	r.POST("/join", func(c *gin.Context) {
		input := JoinInput{}
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
			return
		}

		status := room.JoinRoom(c, input.RoomCode, input.Name)
		if status != http.StatusOK {
			c.JSON(status, gin.H{"message": "Failed to join"})
		}
	})

	r.RunTLS(addr, "cert.pem", "key.pem")
}
