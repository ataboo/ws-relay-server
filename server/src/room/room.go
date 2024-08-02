package room

const RoomLimit = 10
const CodeDigits = 6

type Room struct {
	Code string
	Game *Game
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
