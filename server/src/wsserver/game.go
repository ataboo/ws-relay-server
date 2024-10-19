package wsserver

import (
	"github.com/ataboo/rtc-game-buzzer/src/wsmessage"
)

type Game interface {
	Start() error
	Stop()
	AddPlayer(player *Player) error
	RemovePlayer(id uint16) error
	PlayerCount() int
	Done() chan struct{}
}

type Player struct {
	ID            uint16
	Name          string
	MsgToPlayer   chan *wsmessage.WSMessage
	MsgFromPlayer chan *wsmessage.WSMessage
}
