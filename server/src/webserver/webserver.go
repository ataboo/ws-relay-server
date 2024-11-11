package webserver

import (
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/ataboo/rtc-game-buzzer/src/internal/common"
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
	common.LoadDotEnv()

	addr := os.Getenv("HOSTNAME")
	logLevelStr := os.Getenv("LOG_LEVEL")
	if logLevelStr == "" {
		logLevelStr = "debug"
	}

	fmt.Print("addr", addr)

	logLevel, err := log.ParseLevel(logLevelStr)
	if err != nil {
		log.Errorf("invalid log level: %s", logLevelStr)
		logLevel = log.InfoLevel
	}

	log.SetLevel(logLevel)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	r.Static("static", "./static")

	wsServer = wsserver.NewWsServer(gameFactory)
	wsServer.Start()

	r.GET("/ws", handleWs)

	certPath, err := common.GetAndMakeLocalDir("certs")
	if err != nil {
		log.Fatal(err)
	}

	r.RunTLS(addr, path.Join(certPath, "ws-relay-server.pem.pub"), path.Join(certPath, "ws-relay-server.pem"))
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
