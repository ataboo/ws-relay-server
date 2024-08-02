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

func (u *User) readPump(reqChan chan<- WSReq) {
	defer func() {
		u.conn.Close()
		fmt.Printf("read close: %s\n", u.id)
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
			if err != nil {
				fmt.Print("failed to unmarshal message")
			} else {

				input := WSReq{
					Msg:    msg,
					Sender: u,
				}

				reqChan <- input
			}
		}
	}
}

func (u *User) writePump(leaveChan chan<- *User) {
	ticker := time.NewTicker(PingPeriod)
	defer func() {
		leaveChan <- u
		ticker.Stop()
		fmt.Printf("write close: %s\n", u.id)
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
