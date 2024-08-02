package room

// func TestMakeRoomCode(t *testing.T) {
// 	code := makeRoomCode()

// 	if len(code) != CodeDigits {
// 		t.Errorf("Code unnexpected length %d, %d", len(code), CodeDigits)
// 	}

// 	for _, d := range code {
// 		if d < 'A' || d > 'Z' {
// 			t.Errorf("Digit out of range: '%c'", d)
// 		}
// 	}
// }

// func _setupTestRoomOne(t *testing.T) (rooms *Lobby, client *websocket.Conn, srv *httptest.Server, deferFunc func()) {
// 	roomList := NewLobby()
// 	var roomCode string

// 	gin.SetMode(gin.TestMode)

// 	g := gin.Default()
// 	g.GET("/create", func(ctx *gin.Context) {
// 		roomList.CreateAndStartRoom(ctx, "Player1")
// 	})

// 	g.GET("/join", func(ctx *gin.Context) {
// 		roomList.JoinRoom(ctx, roomCode, "Player2")
// 	})

// 	srv = httptest.NewServer(g.Handler())

// 	url := "ws" + strings.TrimPrefix(srv.URL, "http") + "/create"

// 	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	wrtErr := ws.WriteMessage(websocket.TextMessage, []byte(`{"msg": "foo"}`))
// 	if wrtErr != nil {
// 		t.Error(wrtErr)
// 	}

// 	if len(roomList.rooms) != 1 {
// 		t.Errorf("Unnexpected room length %d", len(roomList.rooms))
// 	}

// 	var room *Room = nil
// 	for _, r := range roomList.rooms {
// 		room = r
// 		break
// 	}

// 	roomCode = room.Code

// 	return roomList, ws, srv, func() {
// 		srv.Close()
// 		ws.Close()
// 	}
// }

// func TestCreateAndStartRoom(t *testing.T) {
// 	roomList, ws, srv, deferFunc := _setupTestRoomOne(t)
// 	defer deferFunc()

// 	wrtErr := ws.WriteMessage(websocket.TextMessage, []byte(`{"msg": "foo"}`))
// 	if wrtErr != nil {
// 		t.Error(wrtErr)
// 	}

// 	if len(roomList.rooms) != 1 {
// 		t.Errorf("Unnexpected room length %d", len(roomList.rooms))
// 	}

// 	var room *Room = nil
// 	for _, r := range roomList.rooms {
// 		room = r
// 		break
// 	}

// 	if room.Game.Host.Name != "Player1" {
// 		t.Errorf("Unnexpected host")
// 	}

// 	u2 := "ws" + strings.TrimPrefix(srv.URL, "http") + "/join"

// 	ws2, _, err := websocket.DefaultDialer.Dial(u2, nil)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	defer ws2.Close()

// 	wrtErr2 := ws2.WriteMessage(websocket.TextMessage, []byte(`{"msg": "foo"}`))
// 	if wrtErr2 != nil {
// 		t.Error(wrtErr)
// 	}

// 	if room.Game.Players[1].Name != "Player2" {
// 		t.Errorf("Unnexpected player 2")
// 	}

// }

// func TestNoDuplicatePlayers(t *testing.T) {
// 	roomList, ws, srv, deferFunc := _setupTestRoomOne(t)
// 	defer deferFunc()

// 	u2 := "ws" + strings.TrimPrefix(srv.URL, "http") + "/join"

// 	ws2, _, err := websocket.DefaultDialer.Dial(u2, nil)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	defer ws2.Close()

// 	wrtErr2 := ws2.WriteMessage(websocket.TextMessage, []byte(`{"msg": "foo"}`))
// 	if wrtErr2 != nil {
// 		t.Error(wrtErr2)
// 	}

// 	var room *Room
// 	for _, r := range roomList.rooms {
// 		room = r
// 		break
// 	}

// 	if room.Game.Players[1].Name != "Player2" {
// 		t.Errorf("Unnexpected player 2")
// 	}
// }
