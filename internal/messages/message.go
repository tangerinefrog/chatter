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
	ReadAt    *time.Time
}
