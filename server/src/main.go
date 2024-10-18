package main

import (
	"github.com/ataboo/rtc-game-buzzer/src/game"
	"github.com/ataboo/rtc-game-buzzer/src/webserver"
)

func main() {
	webserver.Start(game.NewTestGame)
}
