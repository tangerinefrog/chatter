package messages

import "time"

type Message struct {
	ID        int64
	UserID    int32
	ChatID    int32
	Content   string
	CreatedAt time.Time
	ReadAt    *time.Time
}
