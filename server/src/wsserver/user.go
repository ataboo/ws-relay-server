package wsserver

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ataboo/rtc-game-buzzer/src/wsmessage"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type User struct {
	conn        *websocket.Conn
	id          uint16
	name        string
	msgToUser   chan *wsmessage.WSMessage
	msgFromUser chan *wsmessage.WSMessage
	roomId      string
}

func (u *User) Name() string {
	return u.name
}

func (u *User) ID() uint16 {
	return u.id
}

func (u *User) handshake() (joinPayload wsmessage.JoinPayload, err error) {
	joinPayload = wsmessage.JoinPayload{}

	log.Debug("sending welcome message")

	// 1. Server sends welcome to user with their assigned user id.
	welcomePayload := wsmessage.WelcomePayload{UserId: u.id}
	welcomeBytes, err := json.Marshal(welcomePayload)
	if err != nil {
		return joinPayload, err
	}
	welcomeMsg, err := wsmessage.Marshal(wsmessage.CodeWelcome, 0, welcomeBytes)
	if err != nil {
		return joinPayload, err
	}

	// 2. User responds by setting their user name and room code
	u.conn.SetWriteDeadline(time.Now().Add(WriteWait))
	if err := u.conn.WriteMessage(websocket.BinaryMessage, welcomeMsg); err != nil {
		log.Debug("failed to send welcome message")
		return joinPayload, err
	}

	log.Debug("waiting for join message")

	u.conn.SetReadDeadline(time.Now().Add(ReadWait))
	u.conn.SetReadLimit(MaxMessageSize)

	mType, p, err := u.conn.ReadMessage()
	if err != nil {
		log.Debug("failed to read join message")
		return joinPayload, err
	}

	err = wsmessage.ParseMessageWithPayload(mType, p, wsmessage.CodeJoin, &joinPayload)
	if err != nil {
		log.Debug("failed to parse join message")
		return joinPayload, err
	}

	log.Debug("successfully received join message")

	return joinPayload, nil
}

func (u *User) Stop() {
	u.conn.Close()
	if u.msgFromUser != nil {
		close(u.msgFromUser)
		u.msgFromUser = nil
	}
}

func (u *User) readPump() {
	defer func() {
		u.conn.Close()
		fmt.Printf("user read close: %d\n", u.id)
	}()

	u.conn.SetReadLimit(MaxMessageSize)
	u.conn.SetReadDeadline(time.Now().Add(PongWait))
	u.conn.SetPongHandler(func(string) error { u.conn.SetReadDeadline(time.Now().Add(PongWait)); return nil })
	for {
		mType, p, err := u.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("Unexpected close")
			}
			fmt.Printf("client %s read err: %s\n", u.Name(), err.Error())
			return
		}

		if mType == websocket.BinaryMessage {
			msg, err := wsmessage.Unmarshal(p)
			if err == nil {
				if msg.Sender != u.id {
					fmt.Printf("invalid sender")
					return
				}

				u.msgFromUser <- msg
			}
		}
	}
}

func (u *User) writePump(leave chan<- uint16) {
	ticker := time.NewTicker(PingPeriod)
	defer func() {
		leave <- u.ID()
		ticker.Stop()
		fmt.Printf("write close: %d\n", u.id)
	}()

	for {
		select {
		case <-ticker.C:
			u.conn.SetWriteDeadline(time.Now().Add(WriteWait))
			if err := u.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		case res, ok := <-u.msgToUser:
			if !ok {
				u.conn.SetWriteDeadline(time.Now().Add(WriteWait))
				u.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			u.conn.SetWriteDeadline(time.Now().Add(WriteWait))
			if err := u.conn.WriteJSON(res); err != nil {
				return
			}
		}
	}
}
