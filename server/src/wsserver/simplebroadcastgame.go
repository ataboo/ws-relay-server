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
	if g.userNameTaken(player.Name) {
		return fmt.Errorf("player name taken")
	}

	g.players[player.ID] = player

	go func() {
		for {
			msg, ok := <-player.MsgFromPlayer
			if !ok {
				break
			}

			if msg.Code == wsmessage.CodeBroadcast || msg.Code == wsmessage.CodeBroadcastOthers {
				msg.Sender = player.ID
				g.broadcast <- msg
			}
		}
	}()

	g.broadcastPlayerChange()

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

func (g *SimpleBroadcastGame) RemovePlayer(id uint16) error {
	delete(g.players, id)

	if g.PlayerCount() > 0 {
		g.broadcastPlayerChange()
	}

	return nil
}

func (g *SimpleBroadcastGame) Done() chan struct{} {
	return g.doneChan
}

func (g *SimpleBroadcastGame) Start() error {
	if g.running {
		return fmt.Errorf("already running")
	}

	g.running = true

	go func() {
		for {
			msg, ok := <-g.broadcast
			if !ok {
				return
			}

			for _, p := range g.players {
				if msg.Code == wsmessage.CodeBroadcast || p.ID != msg.Sender {
					p.MsgToPlayer <- msg
				}
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

func (g *SimpleBroadcastGame) PlayerCount() int {
	return len(g.players)
}

func (g *SimpleBroadcastGame) broadcastPlayerChange() error {
	pld := wsmessage.PlayerChangePayload{
		Players: make([]wsmessage.PlayerPayload, g.PlayerCount()),
	}

	idx := 0
	for _, p := range g.players {
		pld.Players[idx] = wsmessage.PlayerPayload{Name: p.Name, Id: p.ID}
		idx++
	}

	msg, err := wsmessage.NewWsMessage(wsmessage.CodeBroadcast, wsmessage.ServerSenderId, wsmessage.PldIdPlayerChange, pld)
	if err != nil {
		return err
	}

	g.broadcast <- msg

	return nil
}

func (g *SimpleBroadcastGame) userNameTaken(userName string) bool {
	for _, v := range g.players {
		if v.Name == userName {
			return true
		}
	}

	return false
}
