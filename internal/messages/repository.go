package messages

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tangerinefrog/chatter/internal/crypto"
	"github.com/tangerinefrog/chatter/internal/db"
)

type MessagesRepository struct {
	q      *db.Queries
	cipher *crypto.Cipher
}

const PageSize int32 = 20

func NewRepository(pool *pgxpool.Pool, cipher *crypto.Cipher) *MessagesRepository {
	return &MessagesRepository{
		q:      db.New(pool),
		cipher: cipher,
	}
}

func (r *MessagesRepository) CreateMessage(
	ctx context.Context,
	userID int32,
	chatID int32,
	content string,
) (int64, error) {
	encryptedBytes, err := r.cipher.Encrypt([]byte(content))
	if err != nil {
		return 0, err
	}

	id, err := r.q.CreateMessage(ctx, db.CreateMessageParams{
		ChatID:  chatID,
		UserID:  pgtype.Int4{Int32: userID, Valid: true},
		Content: base64.StdEncoding.EncodeToString(encryptedBytes),
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

		var readAt *time.Time
		if !m.ReadAt.Valid {
			readAt = nil
		} else {
			readAt = &m.ReadAt.Time
		}

		contentBytes, err := base64.StdEncoding.DecodeString(m.Content)
		if err != nil {
			return nil, err
		}

		contentDecrypted, err := r.cipher.Decrypt(contentBytes)
		if err != nil {
			return nil, err
		}
		result[i] = Message{
			ID:        m.ID,
			UserID:    userID,
			ChatID:    chatID,
			Content:   string(contentDecrypted),
			CreatedAt: m.CreatedAt.Time,
			ReadAt:    readAt,
		}
	}

	return result, nil
}

func (r *MessagesRepository) MarkMessagesAsRead(
	ctx context.Context,
	messageID int64,
	chatID int32,
	userID int32,
) error {
	err := r.q.MarkMessagesAsRead(ctx, db.MarkMessagesAsReadParams{
		ID:     messageID,
		ChatID: chatID,
		UserID: pgtype.Int4{Int32: userID, Valid: true},
	})

	if err != nil {
		return err
	}

	return nil
}
