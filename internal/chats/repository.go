package chats

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tangerinefrog/chatter/internal/crypto"
	"github.com/tangerinefrog/chatter/internal/db"
)

type ChatsRepository struct {
	q      *db.Queries
	pool   *pgxpool.Pool
	cipher *crypto.Cipher
}

func NewRepository(p *pgxpool.Pool, cipher *crypto.Cipher) *ChatsRepository {
	return &ChatsRepository{
		q:      db.New(p),
		pool:   p,
		cipher: cipher,
	}
}

func (r *ChatsRepository) CreateChat(ctx context.Context, name string, chatType ChatType, userID uuid.UUID, participantIDs []uuid.UUID) (uuid.UUID, error) {
	if len(participantIDs) < 2 {
		return uuid.Nil, errors.New("chat should have at least 2 participants")
	}

	if len(participantIDs) > 2 && chatType == ChatTypeDirect {
		return uuid.Nil, errors.New("direct chat can have only 2 participants")
	}

	var nameVal pgtype.Text
	if name != "" {
		nameVal = pgtype.Text{String: name}
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback(ctx)

	qtx := r.q.WithTx(tx)

	chatID, err := qtx.CreateChat(ctx, db.CreateChatParams{
		Type:      string(chatType),
		Name:      nameVal,
		CreatedBy: pgtype.UUID{Bytes: userID, Valid: true},
	})

	if err != nil {
		return uuid.Nil, err
	}

	for _, id := range participantIDs {
		err := qtx.CreateChatUser(ctx, db.CreateChatUserParams{ChatID: chatID, UserID: pgtype.UUID{Bytes: id, Valid: true}})
		if err != nil {
			return uuid.Nil, err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return uuid.Nil, err
	}

	return chatID.Bytes, nil
}

func (r *ChatsRepository) IsDirectChatExists(ctx context.Context, userID_1, userID_2 uuid.UUID) (bool, error) {
	existingChatID, err := r.q.GetDirectChatBetweenUsers(ctx, db.GetDirectChatBetweenUsersParams{
		UserID:   pgtype.UUID{Bytes: userID_1, Valid: true},
		UserID_2: pgtype.UUID{Bytes: userID_2, Valid: true},
	})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, errors.New("could not get direct chat between users")
	}
	if existingChatID.Bytes != uuid.Nil {
		return true, nil
	}

	return false, nil
}

func (r *ChatsRepository) ListChatsForUser(ctx context.Context, userID uuid.UUID) ([]Chat, error) {
	chatRows, err := r.q.ListUserChats(ctx, pgtype.UUID{Bytes: userID, Valid: true})
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

		if lastMessage.Content != "" {
			contentBytes, err := base64.StdEncoding.DecodeString(lastMessage.Content)
			if err != nil {
				return nil, err
			}
			contentDecrypted, err := r.cipher.Decrypt(contentBytes)
			if err != nil {
				return nil, err
			}
			lastMessage.Content = string(contentDecrypted)
		}

		chats[i] = Chat{
			ID:                  c.ID.Bytes,
			Type:                c.Type,
			Name:                c.Name.String,
			LastMessage:         lastMessage.Content,
			LastMessageDate:     lastMessage.CreatedAt.UTC(),
			Participants:        participants,
			CreatedAt:           c.CreatedAt.Time,
			UnreadMessagesCount: int32(c.UnreadMessagesCount),
		}
	}

	return chats, nil
}

func (r *ChatsRepository) ListUsersForChat(ctx context.Context, chatID uuid.UUID) ([]*db.User, error) {
	rows, err := r.q.ListChatUsers(ctx, pgtype.UUID{Bytes: chatID, Valid: true})
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
