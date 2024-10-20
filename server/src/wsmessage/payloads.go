package wsmessage

const (
	PldIdWelcome      = 1
	PldIdJoin         = 2
	PldIdPlayerChange = 3
)

type WelcomePayload struct {
	UserId uint16 `json:"user_id"`
}

type JoinPayload struct {
	Name     string `json:"name"`
	RoomCode string `json:"room_code"`
}

type PlayerChangePayload struct {
	Players []PlayerPayload `json:"players"`
}

type PlayerPayload struct {
	Name string `json:"name"`
	Id   uint16 `json:"id"`
}
