package room

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	MaxMessageSize = 2048
	ReadWait       = 3 * time.Second
	WriteWait      = 3 * time.Second
	PongWait       = 10 * time.Second
	PingPeriod     = 5 * time.Second
)

type Lobby struct {
	users     []*User
	rooms     map[string]*Room
	msgChan   chan WSReq
	leaveChan chan *User
}

var randSrc = rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixNano())))

func NewLobby() *Lobby {
	l := &Lobby{
		rooms:     map[string]*Room{},
		users:     []*User{},
		msgChan:   make(chan WSReq),
		leaveChan: make(chan *User),
	}

	return l
}

func (l *Lobby) Start() {
	func() {

		defer func() {
			for _, u := range l.users {
				close(u.writeChan)
			}
		}()

		for {
			msg, ok := <-l.msgChan
			if !ok {
				break
			}

			l.handleMsg(msg)
		}
	}()
}

func (l *Lobby) handleMsg(msg WSReq) {
	fmt.Printf("got message: %d", msg.Msg.Code)
}

func (l *Lobby) Stop() {
	close(l.msgChan)
}

func (l *Lobby) AddUser(conn *websocket.Conn) {
	id, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}

	u := &User{
		id:   id,
		conn: conn,
		name: id.String(),
	}

	l.users = append(l.users, u)

	go u.readPump(l.leaveChan, l.msgChan)
	go u.writePump(l.leaveChan)
}

// func (r *Lobby) JoinRoom(c *gin.Context, roomCode string, name string) int {
// 	room, ok := r.rooms[roomCode]
// 	if !ok {
// 		return http.StatusNotFound
// 	}

// 	if room.Game.Locked {
// 		return http.StatusForbidden
// 	}

// 	for _, p := range room.Game.Players {
// 		if p.Name == name {
// 			return http.StatusForbidden
// 		}
// 	}

// 	conn, err := upgrader.Upgrade(c.Writer, c.Request, c.Request.Header)
// 	if err != nil {
// 		return http.StatusBadRequest
// 	}

// 	player := Player{
// 		IsHost:    true,
// 		Name:      name,
// 		Conn:      conn,
// 		WriteChan: make(chan WSMessage),
// 	}

// 	err = room.Game.AddPlayer(&player)
// 	if err != nil {
// 		return http.StatusBadRequest
// 	}

// 	return http.StatusOK
// }

func makeRoomCode() string {
	code := make([]byte, CodeDigits)
	for d := 0; d < CodeDigits; d++ {
		code[d] = 'A' + byte(randSrc.IntN('Z'-'A'))
	}

	return string(code)
}
