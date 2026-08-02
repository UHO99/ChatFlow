package messages

import "time"

type Type string

const (
	TypeMessage Type = "message"
	TypeJoin    Type = "join"
	TypeLeave   Type = "leave"
)

type IncomingMessage struct {
	Type    Type   `json:"type"`
	Content string `json:"content"`
}

type OutgoingMessage struct {
	Type    Type      `json:"type"`
	RoomID  int64     `json:"room_id"`
	Content string    `json:"content,omitempty"`
	Sender  string    `json:"sender,omitempty"`
	SendAt  time.Time `json:"send_at,omitempty"`
}
