package webserver

import (
	"net/http"
	"os"

	"github.com/ataboo/rtc-game-buzzer/src/wsserver"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
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
	wsServer.Start()

	r.GET("/ws", handleWs)

	r.RunTLS(addr, "cert.pem", "key.pem")
}

func handleWs(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, c.Writer.Header())
	if err != nil {
		log.Warn("failed to upgrade WS")
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed to upgrade"})
		return
	}

	err = wsServer.AddUser(conn)
	if err != nil {
		log.Debugf("failed to add user %s", err.Error())
		conn.Close()
	}
}
