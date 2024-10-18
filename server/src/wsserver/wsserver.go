package wsserver

import (
	"fmt"
	"math/rand/v2"
	"regexp"
	"slices"
	"time"

	"github.com/ataboo/rtc-game-buzzer/src/wsmessage"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type WSServer struct {
	rooms       map[string]*Room
	roomCodes   []string
	gameFactory func() Game
}

type GameFactory func() Game

const MaxRoomCount = 8

var roomCodePattern = regexp.MustCompile(`^[A-Z]{6}$`)

var randSrc = rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixNano())))

func NewWsServer(gameFactory GameFactory) *WSServer {
	l := &WSServer{
		rooms:       map[string]*Room{},
		roomCodes:   []string{},
		gameFactory: gameFactory,
	}

	return l
}

func (w *WSServer) AddUser(conn *websocket.Conn) error {
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

	joinPayload, err := u.handshake()
	if err != nil {
		return err
	}

	if !w.userNameValid(joinPayload.Name) {
		return fmt.Errorf("invalid username")
	}

	u.name = joinPayload.Name

	var room *Room
	if joinPayload.RoomCode == "" {
		if len(w.rooms) >= MaxRoomCount {
			return fmt.Errorf("max room count reached")
		}

		roomCode, err := w.generateRoomCode()
		if err != nil {
			return err
		}

		room = NewRoom(roomCode)
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

		room = oldRoom
	}

	err = room.AddUser(u)
	if err != nil {
		return err
	}

	return nil
}

func (w *WSServer) Stop() {
	for _, r := range w.rooms {
		r.Stop()
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
		w.removeRoom(room)
	}()

	return nil
}

func (w *WSServer) removeRoom(room *Room) {
	delete(w.rooms, room.Code)
	codeIdx := slices.Index(w.roomCodes, room.Code)
	if codeIdx >= 0 {
		w.roomCodes = append(w.roomCodes[:codeIdx], w.roomCodes[codeIdx+1:]...)
	}

	fmt.Printf("Room count: %d\n", len(w.rooms))
}

func (w *WSServer) userNameValid(userName string) bool {
	if len(userName) < 3 || len(userName) > 12 {
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
