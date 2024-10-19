package wsserver

import (
	"fmt"
	"math/rand/v2"
	"regexp"
	"slices"
	"time"

	"github.com/ataboo/rtc-game-buzzer/src/wsmessage"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

const (
	MaxNameLength  = 12
	MinNameLength  = 3
	MaxMessageSize = 2048
	ReadWait       = 3 * time.Second
	WriteWait      = 3 * time.Second
	PongWait       = 10 * time.Second
	PingPeriod     = 1 * time.Second
)

type WSServer struct {
	rooms       map[string]*Room
	roomCodes   []string
	gameFactory func() Game
	users       map[uint16]*User
	leaveChan   chan uint16
}

type GameFactory func() Game

const MaxRoomCount = 8

var roomCodePattern = regexp.MustCompile(`^[A-Z]{6}$`)

var randSrc = rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixNano())))

func NewWsServer(gameFactory GameFactory) *WSServer {
	log.SetLevel(log.DebugLevel)

	s := &WSServer{
		rooms:       map[string]*Room{},
		roomCodes:   []string{},
		gameFactory: gameFactory,
		users:       map[uint16]*User{},
		leaveChan:   make(chan uint16),
	}

	return s
}

func (w *WSServer) AddUser(conn *websocket.Conn) error {
	id := w.getNextUserId()

	u := &User{
		id:          id,
		conn:        conn,
		name:        "",
		msgToUser:   make(chan *wsmessage.WSMessage),
		msgFromUser: make(chan *wsmessage.WSMessage),
		roomId:      "",
	}

	log.Debugf("started handshake with user '%d'", id)

	joinPayload, err := u.handshake()
	if err != nil {
		return err
	}

	if !w.userNameValid(joinPayload.Name) {
		return fmt.Errorf("invalid username")
	}

	u.name = joinPayload.Name

	log.Debugf("user '%d' name '%s'", u.id, u.name)

	var room *Room
	if joinPayload.RoomCode == "" {
		if len(w.rooms) >= MaxRoomCount {
			return fmt.Errorf("max room count reached")
		}

		roomCode, err := w.generateRoomCode()
		if err != nil {
			return err
		}

		log.Debugf("creating new room '%s'", roomCode)

		room = NewRoom(roomCode, w.gameFactory())
		err = w.addRoom(room)
		if err != nil {
			return err
		}
	} else {
		if !w.roomCodeIsValid(joinPayload.RoomCode) {
			return fmt.Errorf("invalid room code")
		}

		oldRoom, ok := w.rooms[joinPayload.RoomCode]
		if !ok {
			return fmt.Errorf("invalid room code")
		}

		log.Debugf("joining room %s", joinPayload.RoomCode)

		room = oldRoom
	}

	u.roomId = room.Code

	p := Player{
		ID:            u.id,
		Name:          u.name,
		MsgToPlayer:   u.msgToUser,
		MsgFromPlayer: u.msgFromUser,
	}
	err = room.AddPlayer(&p)
	if err != nil {
		return err
	}

	go u.readPump()
	go u.writePump(w.leaveChan)

	w.users[u.id] = u

	log.Debugf("successfully added user %d", u.id)

	return nil
}

func (w *WSServer) Start() <-chan struct{} {
	stopChan := make(chan struct{})

	go func() {
		for {
			id, ok := <-w.leaveChan
			if !ok {
				break
			}

			u := w.users[id]
			r := w.rooms[u.roomId]

			if r != nil {
				r.RemovePlayer(id)

				if len(r.Players()) == 0 {
					r.Stop()
				}
			}

			u.Stop()
		}
	}()

	return stopChan
}

func (w *WSServer) getNextUserId() uint16 {
	for i := uint16(1); i <= uint16(1<<16-1); i++ {
		if _, ok := w.users[i]; !ok {
			return i
		}
	}

	panic("failed to get user id")
}

func (w *WSServer) Stop() {
	for _, r := range w.rooms {
		r.Stop()
	}

	for _, u := range w.users {
		u.Stop()
	}
}

func (w *WSServer) addRoom(room *Room) error {
	if len(w.rooms) >= MaxRoomCount {
		return fmt.Errorf("room limit exceeded")
	}

	w.rooms[room.Code] = room
	w.roomCodes = append(w.roomCodes, room.Code)
	roomStopChan := room.Start()

	go func() {
		<-roomStopChan
		fmt.Printf("Stop chan fired!")
		for _, p := range room.Players() {
			w.users[p.ID].Stop()
		}
		w.removeRoom(room)
		fmt.Printf("Room removed!")
	}()

	return nil
}

func (w *WSServer) removeRoom(room *Room) {
	delete(w.rooms, room.Code)
	codeIdx := slices.Index(w.roomCodes, room.Code)
	if codeIdx >= 0 {
		w.roomCodes = append(w.roomCodes[:codeIdx], w.roomCodes[codeIdx+1:]...)
	}
}

func (w *WSServer) userNameValid(userName string) bool {
	if len(userName) < MinNameLength || len(userName) > MaxNameLength {
		return false
	}

	return true
}

func (w *WSServer) roomCodeIsValid(roomCode string) bool {
	return roomCodePattern.Match([]byte(roomCode))
}

func (w *WSServer) generateRoomCode() (string, error) {
	for try := 0; try < 100; try++ {
		code := make([]byte, CodeDigits)
		for d := 0; d < CodeDigits; d++ {
			code[d] = 'A' + byte(randSrc.IntN('Z'-'A'))
		}

		codeStr := string(code)

		if _, ok := w.rooms[codeStr]; !ok {
			return codeStr, nil
		}
	}

	return "", fmt.Errorf("failed to generate room code")
}
