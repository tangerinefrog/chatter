package dto

import "time"

type CreateMessageRequest struct {
	Content string `json:"content" validate:"required,min=1"`
}

type ReadMessageRequest struct {
	MessageID string `json:"message_id" validate:"required"`
}

type CreateMessageResponse struct {
	MessageID string    `json:"message_id"`
	CreatedAt time.Time `json:"created_at"`
}

type ListChatMessagesResponse struct {
	Messages []Message `json:"messages"`
}

type Message struct {
	ID        string     `json:"id"`
	Content   string     `json:"content"`
	FromMe    bool       `json:"from_me"`
	UserID    string     `json:"user_id"`
	CreatedAt time.Time  `json:"created_at"`
	ReadAt    *time.Time `json:"read_at"`
}
