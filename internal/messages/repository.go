package messages

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"time"

	"github.com/google/uuid"
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
	userID uuid.UUID,
	chatID uuid.UUID,
	content string,
) (uuid.UUID, error) {
	encryptedBytes, err := r.cipher.Encrypt([]byte(content))
	if err != nil {
		return uuid.Nil, err
	}

	id, err := r.q.CreateMessage(ctx, db.CreateMessageParams{
		ChatID:  pgtype.UUID{Bytes: chatID, Valid: true},
		UserID:  pgtype.UUID{Bytes: userID, Valid: true},
		Content: base64.StdEncoding.EncodeToString(encryptedBytes),
	})

	if err != nil {
		return uuid.Nil, err
	}

	return id.Bytes, nil
}

func (r *MessagesRepository) ListChatMessages(
	ctx context.Context,
	chatID uuid.UUID,
	pageNumber int32,
) ([]Message, error) {
	offset := PageSize * (pageNumber - 1)

	rows, err := r.q.ListTopNMessages(ctx, db.ListTopNMessagesParams{
		ChatID: pgtype.UUID{Bytes: chatID, Valid: true},
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
		var userID uuid.UUID
		if !m.UserID.Valid {
			userID = uuid.Nil
		} else {
			userID = m.UserID.Bytes
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
			ID:        m.ID.Bytes,
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
	messageID uuid.UUID,
	chatID uuid.UUID,
	userID uuid.UUID,
) error {
	err := r.q.MarkMessagesAsRead(ctx, db.MarkMessagesAsReadParams{
		ID:     pgtype.UUID{Bytes: messageID, Valid: true},
		ChatID: pgtype.UUID{Bytes: chatID, Valid: true},
		UserID: pgtype.UUID{Bytes: userID, Valid: true},
	})

	if err != nil {
		return err
	}

	return nil
}
