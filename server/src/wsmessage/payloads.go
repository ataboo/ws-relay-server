package wsmessage

type WelcomePayload struct {
	UserId uint16 `json:"user_id"`
}

type JoinPayload struct {
	Name     string `json:"name"`
	RoomCode string `json:"room_code"`
}

type RoomUpdatePayload struct {
	RoomCode string          `json:"room_code"`
	Players  []PlayerPayload `json:"players"`
}

type PlayerPayload struct {
	Name string `json:"name"`
	Id   uint16 `json:"id"`
}
