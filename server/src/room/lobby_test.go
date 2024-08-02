package room

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ataboo/rtc-game-buzzer/src/wsmessage"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func TestAddUser(t *testing.T) {
	lobby, srv, deferFunc := _setupTestServer(t)
	url := "ws" + strings.TrimPrefix(srv.URL, "http") + "/ws"
	client, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Error(err)
	}

	defer func() {
		deferFunc()
	}()

	<-time.After(time.Millisecond * 10)

	usr1 := lobby.users[0]

	msg, _ := wsmessage.Marshal(wsmessage.CodeSetName, []byte("Player 1"))
	client.WriteMessage(websocket.BinaryMessage, msg)

	<-time.After(time.Millisecond * 10)

	if usr1.Name() != "Player 1" {
		t.Errorf("unexpected name %s", usr1.name)
	}
	client.Close()
	lobby.Stop()

	<-time.After(time.Millisecond * 10)
}

func _setupTestServer(t *testing.T) (lobby *Lobby, srv *httptest.Server, deferFunc func()) {
	lobby = NewLobby()
	lobby.Start()

	gin.SetMode(gin.TestMode)

	g := gin.Default()
	g.GET("/ws", func(ctx *gin.Context) {
		upgrader := &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		}

		conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, ctx.Request.Header)
		if err != nil {
			t.Error(err)
		}

		lobby.AddUser(conn)
	})

	srv = httptest.NewServer(g.Handler())

	return lobby, srv, func() {
		srv.Close()
	}
}
