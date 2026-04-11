package dto

import "time"

type CreateMessageRequest struct {
	Content string `json:"content" validate:"required,min=1"`
}

type ReadMessageRequest struct {
	MessageID int64 `json:"message_id" validate:"required"`
}

type CreateMessageResponse struct {
	MessageID int64     `json:"message_id"`
	CreatedAt time.Time `json:"created_at"`
}

type ListChatMessagesResponse struct {
	Messages []Message `json:"messages"`
}

type Message struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	FromMe    bool      `json:"from_me"`
	UserID    int32     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	ReadAt    *time.Time `json:"read_at"`
}
