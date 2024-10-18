package game

import (
	"github.com/ataboo/rtc-game-buzzer/src/wsmessage"
	"github.com/ataboo/rtc-game-buzzer/src/wsserver"
	"github.com/google/uuid"
)

var _ wsserver.Game = &TestGame{}

type TestGame struct {
}

func NewTestGame() wsserver.Game {
	return &TestGame{}
}

// AddPlayer implements wsserver.Game.
func (t *TestGame) AddPlayer(player wsserver.Player) error {
	panic("unimplemented")
}

// Done implements wsserver.Game.
func (t *TestGame) Done() chan struct{} {
	panic("unimplemented")
}

// HandleMessage implements wsserver.Game.
func (t *TestGame) HandleMessage(wsmessage.WSMessage) {
	panic("unimplemented")
}

// RemovePlayer implements wsserver.Game.
func (t *TestGame) RemovePlayer(id uuid.UUID) error {
	panic("unimplemented")
}

// Start implements wsserver.Game.
func (t *TestGame) Start() error {
	panic("unimplemented")
}

// Stop implements wsserver.Game.
func (t *TestGame) Stop() {
	panic("unimplemented")
}
