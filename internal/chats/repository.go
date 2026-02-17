package chats

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

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
		CreatedBy: pgtype.Int4{Int32: userID, Valid: true},
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

func (r *ChatsRepository) ListChatsForUser(ctx context.Context, userID int32) ([]Chat, error) {
	chatRows, err := r.q.ListUserChats(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	chats := make([]Chat, len(chatRows))
	for i, c := range chatRows {
		var participants []ChatParticipant
		err := json.Unmarshal(c.ChatUsersJson, &participants)
		if err != nil {
			return nil, err
		}

		var lastMessage struct {
			Content   string    `json:"content"`
			CreatedAt time.Time `json:"created_at"`
		}

		if len(c.LastMessageJson) == 0 {
			c.LastMessageJson = []byte("{}")
		}
		
		err = json.Unmarshal(c.LastMessageJson, &lastMessage)
		if err != nil {
			return nil, err
		}

		chats[i] = Chat{
			ID:              c.ID,
			Type:            c.Type,
			Name:            c.Name.String,
			LastMessage:     lastMessage.Content,
			LastMessageDate: lastMessage.CreatedAt.UTC(),
			Participants:    participants,
			CreatedAt:       c.CreatedAt.Time,
		}
	}

	return chats, nil
}

func (r *ChatsRepository) ListUsersForChat(ctx context.Context, chatID int32) ([]*db.User, error) {
	rows, err := r.q.ListChatUsers(ctx, chatID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	users := make([]*db.User, len(rows))

	for i, u := range rows {
		users[i] = &db.User{
			ID:       u.UserID,
			Username: u.Username,
		}
	}

	return users, nil
}
