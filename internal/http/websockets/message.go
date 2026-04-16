package websockets

import "time"

type Event struct {
	Type      EventType `json:"type"`
	Content   string    `json:"content"`
	ChatID    int32     `json:"chat_id,omitempty"`
	SenderID  int32     `json:"sender_id,omitempty"`
	MessageID int64     `json:"message_id,omitempty"`
	FromMe    bool      `json:"from_me,omitempty"`
	ReadAt    time.Time `json:"read_at,omitempty"`
}

type EventType string

const (
	EventSendMessage EventType = "send_message"
	EventNewMessage  EventType = "new_message"
	EventReadMessage EventType = "read_message"
)
