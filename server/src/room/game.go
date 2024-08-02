package room

type Game struct {
	Host    *Player
	Players []*Player
	Locked  bool
}

type Player struct {
	IsHost bool
	User   *User
}

func (g *Game) Start() {

}

func (g *Game) Stop() {

}

func (g *Game) AddPlayer(player *Player) error {
	// for _, p := range g.Players {
	// 	if p.Name == player.Name {
	// 		return fmt.Errorf("name duplicate")
	// 	}
	// }

	// if len(g.Players) == 0 {
	// 	g.Host = player
	// }
	// g.Players = append(g.Players, player)

	// go player.readPump(g.leaveChan, g.msgChan)
	// go player.writePump(g.leaveChan)

	// return nil

	return nil
}

// func NewWSClient(conn *websocket.Conn, id int, name string) *WSClient {
// 	return &WSClient{
// 		conn:      conn,
// 		writeChan: make(chan *msg.WSResponse),
// 		ClientID:  id,
// 		Name:      name,
// 	}
// }

// func (c *WSClient) Start(leaveChan chan<- *WSClient, reqChan chan<- *msg.WSRequest) {
// 	go c.readPump(leaveChan, reqChan)
// 	go c.writePump(leaveChan)
// }

// func (c *WSClient) writePump(leaveChan chan<- *WSClient) {
// 	ticker := time.NewTicker(PingPeriod)
// 	defer func() {
// 		ticker.Stop()
// 		c.conn.Close()
// 	}()

// 	for {
// 		select {
// 		case <-ticker.C:
// 			c.conn.SetWriteDeadline(time.Now().Add(WriteWait))
// 			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
// 				return
// 			}
// 		case res, ok := <-c.writeChan:
// 			if !ok {
// 				c.conn.SetWriteDeadline(time.Now().Add(WriteWait))
// 				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
// 				return
// 			}

// 			c.conn.SetWriteDeadline(time.Now().Add(WriteWait))
// 			if err := c.conn.WriteJSON(res); err != nil {
// 				return
// 			}
// 		}
// 	}
// }

// func (c *WSClient) WriteResponse(res *msg.WSResponse) bool {
// 	select {
// 	case c.writeChan <- res:
// 		return true
// 	default:
// 		close(c.writeChan)
// 		return false
// 	}
// }

// func WriteResponse(conn *websocket.Conn, res msg.WSResponse) error {
// 	conn.SetWriteDeadline(time.Now().Add(WriteWait))
// 	return conn.WriteJSON(res)
// }
