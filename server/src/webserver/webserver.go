package webserver

import (
	"net/http"
	"os"

	"github.com/ataboo/rtc-game-buzzer/src/wsserver"
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

var wsServer *wsserver.WSServer

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func Start(gameFactory wsserver.GameFactory) {
	addr := os.Getenv("HOSTNAME")

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	wsServer = wsserver.NewWsServer(gameFactory)

	r.POST("/ws", handleWs)

	r.RunTLS(addr, "cert.pem", "key.pem")
}

func handleWs(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, c.Request.Header)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed to upgrade"})
		return
	}

	err = wsServer.AddUser(conn)
	if err != nil {
		conn.Close()
	}
}
