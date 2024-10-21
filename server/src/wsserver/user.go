package wsserver

import (
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
		log.Debugf("user read close: %d\n", u.id)
	}()

	u.conn.SetReadLimit(MaxMessageSize)
	u.conn.SetReadDeadline(time.Now().Add(PongWait))
	u.conn.SetPongHandler(func(string) error { u.conn.SetReadDeadline(time.Now().Add(PongWait)); return nil })
	for {
		mType, p, err := u.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Debugf("Unexpected close")
			}
			log.Debugf("client %s read err: %s\n", u.Name(), err.Error())
			return
		}

		if mType == websocket.BinaryMessage {
			msg, err := wsmessage.Unmarshal(p)
			if err != nil {
				log.Debugf("failed to parse msg, %s", err.Error())
				return
			}

			msg.Sender = u.id
			u.msgFromUser <- msg
		}
	}
}

func (u *User) writePump(leave chan<- uint16) {
	ticker := time.NewTicker(PingPeriod)
	defer func() {
		leave <- u.ID()
		ticker.Stop()
		log.Debugf("write close: %d\n", u.id)
	}()

	for {
		select {
		case <-ticker.C:
			u.conn.SetWriteDeadline(time.Now().Add(WriteWait))
			if err := u.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		case msg, ok := <-u.msgToUser:
			if !ok {
				u.conn.SetWriteDeadline(time.Now().Add(WriteWait))
				u.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			msgBytes, err := wsmessage.Marshal(msg)
			if err != nil {
				log.Warn("failed to marshal outgoing message")
				return
			}

			u.conn.SetWriteDeadline(time.Now().Add(WriteWait))
			if err := u.conn.WriteMessage(websocket.BinaryMessage, msgBytes); err != nil {
				return
			}

			log.Debugf("%d sent to %d", msg.PayloadId, u.id)
		}
	}
}
