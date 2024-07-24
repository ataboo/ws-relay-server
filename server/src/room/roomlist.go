package room

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const RoomLimit = 10

var rooms = map[string]*Room{}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func CreateRoom(c *gin.Context, name string) int {

	if len(rooms) >= RoomLimit {
		return http.StatusForbidden
	}

	hostConn, err := upgrader.Upgrade(c.Writer, c.Request, c.Request.Header)
	if err != nil {
		return http.StatusBadRequest
	}

	hostPlayer := Player{
		IsHost:    true,
		Name:      name,
		Conn:      hostConn,
		WriteChan: make(chan WSMessage),
	}

	roomCode := makeRoomCode()

	room := Room{
		Code: roomCode,
		Game: &Game{
			Players: []*Player{},
			Host:    nil,
			Locked:  false,
		},
	}
	rooms[roomCode] = &room

	room.Game.AddPlayer(&hostPlayer)

	room.Game.Start()

	return http.StatusOK
}

func JoinRoom(c *gin.Context, roomCode string, name string) int {
	room, ok := rooms[roomCode]
	if !ok {
		return http.StatusNotFound
	}

	if room.Game.Locked {
		return http.StatusForbidden
	}

	for _, p := range room.Game.Players {
		if p.Name == name {
			return http.StatusForbidden
		}
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, c.Request.Header)
	if err != nil {
		return http.StatusBadRequest
	}

	player := Player{
		IsHost:    true,
		Name:      name,
		Conn:      conn,
		WriteChan: make(chan WSMessage),
	}

	err = room.Game.AddPlayer(&player)
	if err != nil {
		return http.StatusBadRequest
	}

	return http.StatusOK
}

func makeRoomCode() string {
	return "ABCDEF"
}
