package websockets

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	Type      EventType   `json:"type"`
	Content   string      `json:"content"`
	ChatID    uuid.UUID   `json:"chat_id,omitempty"`
	SenderID  uuid.UUID   `json:"sender_id,omitempty"`
	MessageID uuid.UUID   `json:"message_id,omitempty"`
	FileIDs   []uuid.UUID `json:"file_ids,omitempty"`
	Files     []FileInfo  `json:"files,omitempty"`
	FromMe    bool        `json:"from_me,omitempty"`
	ReadAt    time.Time   `json:"read_at,omitempty"`
}

type FileInfo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	MimeType  string `json:"mime_type"`
	SizeBytes int64  `json:"size_bytes"`
}

type EventType string

const (
	EventSendMessage EventType = "send_message"
	EventNewMessage  EventType = "new_message"
	EventReadMessage EventType = "read_message"
)
