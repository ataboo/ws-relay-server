package wsserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ataboo/rtc-game-buzzer/src/wsmessage"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

func TestHandshakeNewRoom(t *testing.T) {
	log.SetLevel(log.ErrorLevel)

	server, srv, deferFunc := _setupTestServer(t)
	url := "ws" + strings.TrimPrefix(srv.URL, "http") + "/ws"
	client, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Error(err)
	}

	defer func() {
		deferFunc()
	}()

	<-time.After(time.Millisecond * 10)

	mType, p, err := client.ReadMessage()
	if err != nil {
		t.Error(err)
	}

	if mType != websocket.BinaryMessage {
		t.Error("unexpected ws message type")
	}

	welcomeMsg, err := wsmessage.Unmarshal(p)
	if err != nil {
		t.Error(err)
	}

	welcomePL := wsmessage.WelcomePayload{}
	err = json.Unmarshal(welcomeMsg.RawPayload, &welcomePL)
	if err != nil {
		t.Error(err)
	}

	if len(server.rooms) != 0 {
		t.Error("shouldn't be any rooms yet")
	}

	msg, _ := wsmessage.NewWsMessage(wsmessage.CodeJoin, 23, wsmessage.JoinPayload{Name: "Player 1", RoomCode: ""})
	msgBytes, _ := wsmessage.Marshal(msg)

	client.WriteMessage(websocket.BinaryMessage, msgBytes)

	<-time.After(time.Millisecond * 10)

	//TODO get room update message
	// , err := client.ReadMessage()

	room := server.rooms[server.roomCodes[0]]
	user1 := server.users[welcomePL.UserId]

	if !server.roomCodeIsValid(room.Code) {
		t.Error("unexpected room code")
	}

	if user1.Name() != "Player 1" {
		t.Errorf("unexpected name %s", user1.Name())
	}
	client.Close()
	server.Stop()

	<-time.After(time.Millisecond * 10)
}

func TestHandshakeExistingRoom(t *testing.T) {
	log.SetLevel(log.ErrorLevel)

	server, srv, deferFunc := _setupTestServer(t)

	err := server.addRoom(NewRoom("ABCDEF", server.gameFactory()))
	if err != nil {
		t.Error(err)
	}

	url := "ws" + strings.TrimPrefix(srv.URL, "http") + "/ws"
	client, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Error(err)
	}

	defer func() {
		deferFunc()
	}()

	<-time.After(time.Millisecond * 10)

	mType, p, err := client.ReadMessage()
	if err != nil {
		t.Error(err)
	}

	if mType != websocket.BinaryMessage {
		t.Error("unexpected ws message type")
	}

	welcomeMsg, err := wsmessage.Unmarshal(p)
	if err != nil {
		t.Error(err)
	}

	welcomePL := wsmessage.WelcomePayload{}
	err = json.Unmarshal(welcomeMsg.RawPayload, &welcomePL)
	if err != nil {
		t.Error(err)
	}

	msg, _ := wsmessage.NewWsMessage(wsmessage.CodeJoin, 23, wsmessage.JoinPayload{Name: "Player 1", RoomCode: "ABCDEF"})
	msgBytes, _ := wsmessage.Marshal(msg)

	client.WriteMessage(websocket.BinaryMessage, msgBytes)

	<-time.After(time.Millisecond * 10)

	//TODO get room update message
	// , err := client.ReadMessage()

	room := server.rooms[server.roomCodes[0]]
	if room.Code != "ABCDEF" {
		t.Error("unexpected room code")
	}

	user1 := server.users[welcomePL.UserId]

	if user1.Name() != "Player 1" {
		t.Errorf("unexpected name %s", user1.Name())
	}
	client.Close()
	<-time.After(time.Millisecond * 1000)

	fmt.Printf("Hit expect\n")
	if len(server.rooms) > 0 {
		t.Error("expected room to be cleaned up")
	}

	server.Stop()

	<-time.After(time.Millisecond * 10)
}

func _setupTestServer(t *testing.T) (server *WSServer, srv *httptest.Server, deferFunc func()) {
	server = NewWsServer(NewSimpleBroadcastGame)
	server.Start()

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

		server.AddUser(conn)
	})

	srv = httptest.NewServer(g.Handler())

	return server, srv, func() {
		srv.Close()
	}
}
