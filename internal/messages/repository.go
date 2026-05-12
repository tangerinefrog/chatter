package messages

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tangerinefrog/chatter/internal/db"
)

type MessagesRepository struct {
	q             *db.Queries
	encryptionKey []byte
}

const PageSize int32 = 20

func (r *MessagesRepository) deriveChatKey(chatID int32) []byte {
	h := sha256.New()
	h.Write(r.encryptionKey)
	h.Write([]byte(fmt.Sprintf("%d", chatID)))
	return h.Sum(nil)
}

func NewRepository(pool *pgxpool.Pool, encryptionKey []byte) *MessagesRepository {
	return &MessagesRepository{
		q:             db.New(pool),
		encryptionKey: encryptionKey,
	}
}

const encryptionKeySize = 32

func EncryptMessageContent(key []byte, plaintext string) (string, error) {
	if len(key) != encryptionKeySize {
		return "", fmt.Errorf("encryption key must be %d bytes", encryptionKeySize)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create AEAD: %w", err)
	}

	nonce := make([]byte, aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := aead.Seal(nil, nonce, []byte(plaintext), nil)
	payload := append(nonce, ciphertext...)

	return base64.StdEncoding.EncodeToString(payload), nil
}

func DecryptMessageContent(key []byte, encoded string) (string, error) {
	if len(key) != encryptionKeySize {
		return "", fmt.Errorf("encryption key must be %d bytes", encryptionKeySize)
	}

	payload, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("failed to decode ciphertext: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create AEAD: %w", err)
	}

	nonceSize := aead.NonceSize()
	if len(payload) < nonceSize {
		return "", errors.New("invalid ciphertext payload")
	}

	nonce := payload[:nonceSize]
	ciphertext := payload[nonceSize:]

	plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt message: %w", err)
	}

	return string(plaintext), nil
}

func (r *MessagesRepository) CreateMessage(
	ctx context.Context,
	userID int32,
	chatID int32,
	content string,
) (int64, error) {
	chatKey := r.deriveChatKey(chatID)

	encryptedContent, err := EncryptMessageContent(chatKey, content)
	if err != nil {
		return 0, fmt.Errorf("failed to encrypt message: %w", err)
	}

	id, err := r.q.CreateMessage(ctx, db.CreateMessageParams{
		ChatID:  chatID,
		UserID:  pgtype.Int4{Int32: userID, Valid: true},
		Content: encryptedContent,
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

	chatKey := r.deriveChatKey(chatID)

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

		decryptedContent, err := DecryptMessageContent(chatKey, m.Content)
		if err != nil {
			decryptedContent = m.Content
		}

		result[i] = Message{
			ID:        m.ID,
			UserID:    userID,
			ChatID:    chatID,
			Content:   decryptedContent,
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
