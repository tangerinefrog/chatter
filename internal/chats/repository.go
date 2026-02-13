package chats

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tangerinefrog/chatter/internal/db"
)

type ChatsRepository struct {
	q    *db.Queries
	pool *pgxpool.Pool
}

func NewRepository(p *pgxpool.Pool) *ChatsRepository {
	return &ChatsRepository{
		q:    db.New(p),
		pool: p,
	}
}

func (r *ChatsRepository) Create(ctx context.Context, name string, chatType ChatType, userID int32, participantIDs []int32) (int32, error) {
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

func (r *ChatsRepository) IsDirectChatExists(ctx context.Context, userID_1, userID_2 int32) (bool, error) {
	existingChatID, err := r.q.GetDirectChatBetweenUsers(ctx, db.GetDirectChatBetweenUsersParams{
		UserID:   userID_1,
		UserID_2: userID_2,
	})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, errors.New("could not get direct chat between users")
	}
	if existingChatID != 0 {
		return true, nil
	}

	return false, nil
}
