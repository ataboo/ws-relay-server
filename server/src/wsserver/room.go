package wsserver

import (
	"fmt"
)

const (
	CodeDigits = 6
)

type Room struct {
	Code string
	game Game
}

func NewRoom(roomCode string, game Game) *Room {
	room := &Room{
		Code: roomCode,
		game: game,
	}

	return room
}

func (r Room) Start() chan struct{} {
	r.game.Start()

	return r.game.Done()
}

func (r *Room) Stop() {
	r.game.Stop()
}

func (r *Room) AddPlayer(p *Player) error {
	if r.userNameTaken(p.Name) {
		return fmt.Errorf("username taken")
	}

	r.game.AddPlayer(p)

	return nil
}

func (r *Room) userNameTaken(userName string) bool {
	for _, v := range r.game.Players() {
		if v.Name == userName {
			return true
		}
	}

	return false
}

func (r *Room) RemovePlayer(id uint16) error {
	return r.game.RemovePlayer(id)
}

func (r *Room) Players() []*Player {
	return r.game.Players()
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
