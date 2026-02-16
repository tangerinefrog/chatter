package chats

import (
	"time"
)

type Chat struct {
	ID              int32
	Type            string
	Name            string
	Participants    []ChatParticipant
	LastMessage     string
	LastMessageDate time.Time
	CreatedBy       int32
	CreatedAt       time.Time
}

type ChatParticipant struct {
	ID       int32  `json:"id"`
	Username string `json:"username"`
}
