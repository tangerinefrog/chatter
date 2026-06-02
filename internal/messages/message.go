package messages

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	ChatID    uuid.UUID
	Content   string
	CreatedAt time.Time
	Files     []MessageFile
	ReadAt    *time.Time
}

type MessageFile struct {
	ID        uuid.UUID `json:"id"`
	FileKey   string    `json:"file_key"`
	FileName  string    `json:"file_name"`
	MimeType  string    `json:"mime_type"`
	SizeBytes int64     `json:"size_bytes"`
}
