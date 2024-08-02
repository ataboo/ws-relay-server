package room

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"time"

	"github.com/ataboo/rtc-game-buzzer/src/wsmessage"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	MaxMessageSize = 2048
	ReadWait       = 3 * time.Second
	WriteWait      = 3 * time.Second
	PongWait       = 10 * time.Second
	PingPeriod     = 1 * time.Second
	MaxNameLength  = 24
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
	go func() {
		defer func() {
			for _, u := range l.users {
				if u.writeChan != nil {
					close(u.writeChan)
					u.writeChan = nil
				}
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

	go func() {
		for {
			u, ok := <-l.leaveChan
			if !ok {
				break
			}

			idx := slices.IndexFunc(l.users, func(usr *User) bool {
				return u.ID() == usr.ID()
			})

			if idx < 0 {
				fmt.Printf("failed to find user to remove\n")
			} else {
				l.users = append(l.users[:idx], l.users[idx+1:]...)

				if u.writeChan != nil {
					close(u.writeChan)
					u.writeChan = nil
				}
			}
		}
	}()
}

func (l *Lobby) handleMsg(req WSReq) {
	switch req.Msg.Code {
	case wsmessage.CodeSetName:
		l.handleSetName(req)
	default:
		fmt.Printf("unsupported code %d\n", req.Msg.Code)
	}
}

func (l *Lobby) handleSetName(req WSReq) {
	name := string(req.Msg.RawPayload)
	if len(name) == 0 || len(name) > MaxNameLength {
		fmt.Print("invalid name\n")
		return
	}

	req.Sender.name = name

	// broadcast name change
}

func (l *Lobby) Stop() {
	fmt.Print("stop called\n")
	if l.msgChan != nil {
		close(l.msgChan)
		l.msgChan = nil
	}
}

func (l *Lobby) AddUser(conn *websocket.Conn) {
	id, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}

	u := &User{
		id:        id,
		conn:      conn,
		name:      id.String(),
		writeChan: make(chan wsmessage.WSMessage),
	}

	l.users = append(l.users, u)

	go u.readPump(l.msgChan)
	go u.writePump(l.leaveChan)

	fmt.Printf("added user: %s\n", id)
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
