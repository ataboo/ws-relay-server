package room

import (
	"fmt"
	"time"

	"github.com/ataboo/rtc-game-buzzer/src/wsmessage"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type User struct {
	conn      *websocket.Conn
	id        uuid.UUID
	name      string
	writeChan chan wsmessage.WSMessage
}

type WSReq struct {
	Msg    *wsmessage.WSMessage
	Sender *User
}

func (u *User) Name() string {
	return u.name
}

func (u *User) ID() uuid.UUID {
	return u.id
}

func (u *User) readPump(leaveChan chan<- *User, reqChan chan<- WSReq) {
	defer func() {
		leaveChan <- u
		u.conn.Close()
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
			fmt.Printf("Client %s Read err: %s", u.Name(), err.Error())
			return
		}

		if mType == websocket.BinaryMessage {
			msg, err := wsmessage.Unmarshal(p)
			if err != nil {
				fmt.Print("failed to unmarshal message")
			}

			input := WSReq{
				Msg:    msg,
				Sender: u,
			}

			reqChan <- input
		}
	}
}

func (u *User) writePump(leaveChan chan<- *User) {
	ticker := time.NewTicker(PingPeriod)
	defer func() {
		leaveChan <- u
		ticker.Stop()
		u.conn.Close()
	}()

	for {
		select {
		case <-ticker.C:
			u.conn.SetWriteDeadline(time.Now().Add(WriteWait))
			if err := u.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		case res, ok := <-u.writeChan:
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
