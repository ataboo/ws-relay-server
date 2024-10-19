package main

import (
	"github.com/ataboo/rtc-game-buzzer/src/webserver"
	"github.com/ataboo/rtc-game-buzzer/src/wsserver"
)

func main() {
	webserver.Start(wsserver.NewSimpleBroadcastGame)
}
