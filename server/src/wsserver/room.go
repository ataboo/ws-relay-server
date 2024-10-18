package wsserver

import (
	"fmt"
	"sync"
	"time"
)

const (
	MaxMessageSize = 2048
	ReadWait       = 3 * time.Second
	WriteWait      = 3 * time.Second
	PongWait       = 10 * time.Second
	PingPeriod     = 1 * time.Second
	MaxNameLength  = 24
	CodeDigits     = 6
)

type Room struct {
	Code string
	Game *Game

	users     map[string]*User
	msgChan   chan WSReq
	leaveChan chan *User
	lock      sync.Mutex
}

func NewRoom(roomCode string) *Room {
	room := &Room{
		users:     map[string]*User{},
		msgChan:   make(chan WSReq),
		leaveChan: make(chan *User),
		lock:      sync.Mutex{},
		Code:      roomCode,
		Game:      nil,
	}

	return room
}

func (r *Room) Start() <-chan struct{} {
	stopChan := make(chan struct{})

	go func() {
		defer func() {
			for _, u := range r.users {
				if u.writeChan != nil {
					close(u.writeChan)
					u.writeChan = nil
				}
			}

			stopChan <- struct{}{}
		}()

		for {
			msg, ok := <-r.msgChan
			if !ok {
				break
			}

			r.handleMsg(msg)
		}
	}()

	go func() {
		for {
			u, ok := <-r.leaveChan
			if !ok {
				break
			}

			r.removeUser(u)

			if len(r.users) == 0 {
				r.Stop()
				break
			}
		}
	}()

	return stopChan
}

func (r *Room) Stop() {
	fmt.Print("stop called\n")
	if r.msgChan != nil {
		close(r.msgChan)
		r.msgChan = nil
	}

	// if r.leaveChan != nil {
	// 	close(r.leaveChan)
	// 	r.leaveChan = nil
	// }
}

func (r *Room) AddUser(u *User) error {
	if r.userNameTaken(u.name) {
		return fmt.Errorf("username taken")
	}

	r.users[u.id.String()] = u
	go u.readPump(r.msgChan)
	go u.writePump(r.leaveChan)

	// broadcast added user

	return nil
}

func (r *Room) handleMsg(req WSReq) {
	switch req.Msg.Code {
	// case wsmessage.CodeJoin:
	// 	l.handleSetName(req)
	default:
		fmt.Printf("unsupported code %d\n", req.Msg.Code)
	}
}

func (r *Room) userNameTaken(userName string) bool {
	for _, v := range r.users {
		if v.name == userName {
			return true
		}
	}

	return false
}

func (r *Room) removeUser(u *User) error {
	if u.writeChan != nil {
		close(u.writeChan)
		u.writeChan = nil
	}

	delete(r.users, u.id.String())

	return nil
}

// func (l *Lobby) AddNewRoom() *Room {
// 	var roomCode string
// 	for {
// 		roomCode = makeRoomCode()
// 		if _, ok := l.rooms[roomCode]; !ok {
// 			break
// 		}
// 	}

// 	room := Room{
// 		Code: roomCode,
// 		Game: &Game{
// 			Players: []*Player{},
// 			Host:    nil,
// 			Locked:  false,
// 		},
// 	}

// 	l.rooms[roomCode] = &room

// 	return &room
// }
