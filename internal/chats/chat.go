package chats

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ID                  uuid.UUID
	Type                string
	Name                string
	Participants        []ChatParticipant
	LastMessage         string
	LastMessageDate     time.Time
	CreatedBy           uuid.UUID
	CreatedAt           time.Time
	UnreadMessagesCount int32
}

type ChatParticipant struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}
