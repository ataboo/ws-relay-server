package server

import (
	"net/http"
	"os"

	"github.com/ataboo/rtc-game-buzzer/src/room"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type JoinInput struct {
	RoomCode string
	Name     string
}

type HostInput struct {
	Name string
}

var hostAuthKey = os.Getenv("HOST_AUTH_KEY")

var roomList *room.Lobby

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func Start() {
	addr := os.Getenv("HOSTNAME")

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	roomList = room.NewLobby()

	r.POST("/ws", handleWs)

	r.RunTLS(addr, "cert.pem", "key.pem")
}

func handleWs(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, c.Request.Header)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed to upgrade"})
		return
	}

	roomList.AddUser(conn)
}
