package dto

import "time"

type CreateChatRequest struct {
	IsDirect             bool     `json:"is_direct" validate:"required"`
	Name                 string   `json:"name" validate:"max=64"`
	ParticipantUsernames []string `json:"participant_usernames" validate:"required"`
}

type NewChatResponse struct {
	ID string `json:"id"`
}

type ListChatsForUserResponse struct {
	Chats []Chat `json:"chats"`
}

type Chat struct {
	ID                  string            `json:"id"`
	Type                string            `json:"type"`
	Name                string            `json:"name"`
	Participants        []ChatParticipant `json:"participants"`
	LastMessage         string            `json:"last_message"`
	LastMessageDate     *time.Time        `json:"last_message_date"`
	UnreadMessagesCount int32             `json:"unread_messages_count"`
}

type ChatParticipant struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}
