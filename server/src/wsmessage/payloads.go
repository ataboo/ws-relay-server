package wsmessage

type WelcomePayload struct {
	UserId string `json:"user_id"`
}

type JoinPayload struct {
	Name     string `json:"name"`
	GameType uint16 `json:"game_type"`
	RoomCode string `json:"room_code"`
}

type RoomUpdatePayload struct {
	RoomCode string          `json:"room_code"`
	Players  []PlayerPayload `json:"players"`
}

type PlayerPayload struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}
