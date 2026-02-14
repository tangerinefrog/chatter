package messages

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tangerinefrog/chatter/internal/db"
)

type MessagesRepository struct {
	q *db.Queries
}

const PageSize int32 = 100

func NewRepository(pool *pgxpool.Pool) *MessagesRepository {
	return &MessagesRepository{
		q: db.New(pool),
	}
}

func (r *MessagesRepository) CreateMessage(
	ctx context.Context,
	userID int32,
	chatID int32,
	content string,
) (int64, error) {
	id, err := r.q.CreateMessage(ctx, db.CreateMessageParams{
		ChatID:  chatID,
		UserID:  pgtype.Int4{Int32: userID, Valid: true},
		Content: content,
	})

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *MessagesRepository) ListChatMessages(
	ctx context.Context,
	chatID int32,
	pageNumber int32,
) ([]Message, error) {
	offset := PageSize * (pageNumber - 1)

	rows, err := r.q.ListTopNMessages(ctx, db.ListTopNMessagesParams{
		ChatID: chatID,
		Limit:  PageSize,
		Offset: offset,
	})

	fmt.Printf("MESSAGES FROM DB: %v\n", rows)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	result := make([]Message, len(rows))
	for i, m := range rows {
		var userID int32
		if !m.UserID.Valid {
			userID = -1
		} else {
			userID = m.UserID.Int32
		}

		result[i] = Message{
			ID:        m.ID,
			UserID:    userID,
			ChatID:    chatID,
			Content:   m.Content,
			CreatedAt: m.CreatedAt.Time,
		}
	}

	return result, nil
}
