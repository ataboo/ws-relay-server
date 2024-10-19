package wsserver

import (
	"fmt"

	"github.com/ataboo/rtc-game-buzzer/src/wsmessage"
)

var _ Game = &SimpleBroadcastGame{}

type SimpleBroadcastGame struct {
	players   map[uint16]*Player
	doneChan  chan struct{}
	broadcast chan *wsmessage.WSMessage
	running   bool
	stopped   bool
}

func NewSimpleBroadcastGame() Game {
	return &SimpleBroadcastGame{
		players:   map[uint16]*Player{},
		doneChan:  make(chan struct{}),
		broadcast: make(chan *wsmessage.WSMessage),
		running:   false,
	}
}

func (g *SimpleBroadcastGame) AddPlayer(player *Player) error {
	for _, p := range g.players {
		if p.Name == player.Name {
			return fmt.Errorf("player name taken")
		}
	}

	g.players[player.ID] = player

	go func() {
		for {
			msg := <-player.MsgFromPlayer
			for _, p := range g.players {
				if p.ID != player.ID {
					p.MsgToPlayer <- msg
				}
			}
		}
	}()

	return nil
}

func (g *SimpleBroadcastGame) Players() []*Player {
	playerSlice := make([]*Player, len(g.players))

	idx := 0
	for _, v := range g.players {
		playerSlice[idx] = v
		idx++
	}

	return playerSlice
}

func (g *SimpleBroadcastGame) Done() chan struct{} {
	return g.doneChan
}

func (g *SimpleBroadcastGame) RemovePlayer(id uint16) error {
	delete(g.players, id)

	return nil
}

func (g *SimpleBroadcastGame) Start() error {
	if g.running {
		return fmt.Errorf("already running")
	}

	g.running = true

	go func() {
		for {
			msg := <-g.broadcast
			for _, p := range g.players {
				p.MsgToPlayer <- msg
			}
		}
	}()

	return nil
}

func (g *SimpleBroadcastGame) Stop() {
	if g.stopped {
		return
	}

	g.stopped = true
	close(g.broadcast)
	g.doneChan <- struct{}{}
}
