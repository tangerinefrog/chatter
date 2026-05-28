package websockets

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	Type      EventType `json:"type"`
	Content   string    `json:"content"`
	ChatID    uuid.UUID `json:"chat_id,omitempty"`
	SenderID  uuid.UUID `json:"sender_id,omitempty"`
	MessageID uuid.UUID `json:"message_id,omitempty"`
	FileID    uuid.UUID `json:"file_id,omitempty"`
	FromMe    bool      `json:"from_me,omitempty"`
	ReadAt    time.Time `json:"read_at,omitempty"`
}

type EventType string

const (
	EventSendMessage EventType = "send_message"
	EventNewMessage  EventType = "new_message"
	EventReadMessage EventType = "read_message"
)
