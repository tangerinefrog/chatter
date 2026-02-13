package chats

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tangerinefrog/chatter/internal/db"
)

type ChatRepository struct {
	q    *db.Queries
	pool *pgxpool.Pool
}

func NewRepository(p *pgxpool.Pool) *ChatRepository {
	return &ChatRepository{
		q:    db.New(p),
		pool: p,
	}
}

func (r *ChatRepository) Create(ctx context.Context, name string, chatType ChatType, userID int32, participantIDs []int32) (int32, error) {
	if len(participantIDs) < 2 {
		return 0, errors.New("chat should have at least 2 participants")
	}

	if len(participantIDs) > 2 && chatType == ChatTypeDirect {
		return 0, errors.New("direct chat can have only 2 participants")
	}

	var nameVal pgtype.Text
	if name != "" {
		nameVal = pgtype.Text{String: name}
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	qtx := r.q.WithTx(tx)

	chatID, err := qtx.CreateChat(ctx, db.CreateChatParams{
		Type:      string(chatType),
		Name:      nameVal,
		CreatedBy: pgtype.Int4{Int32: userID},
	})

	if err != nil {
		return 0, err
	}

	for _, id := range participantIDs {
		err := qtx.CreateChatUser(ctx, db.CreateChatUserParams{ChatID: chatID, UserID: id})
		if err != nil {
			return 0, err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return 0, err
	}

	return chatID, nil
}
