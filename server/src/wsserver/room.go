package wsserver

import "fmt"

const (
	CodeDigits = 6
)

type Room struct {
	Code string
	Game Game
}

func NewRoom(roomCode string, game Game) *Room {
	room := &Room{
		Code: roomCode,
		Game: game,
	}

	return room
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
