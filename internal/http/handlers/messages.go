package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/tangerinefrog/chatter/internal/http/dto"
	"github.com/tangerinefrog/chatter/internal/messages"
	"go.uber.org/zap"
)

type messagesHandler struct {
	messagesRepo *messages.MessagesRepository
	logger       *zap.Logger
}

func NewMessagesHandler(messagesRepo *messages.MessagesRepository, logger *zap.Logger) *messagesHandler {
	return &messagesHandler{
		messagesRepo: messagesRepo,
		logger:       logger,
	}
}

func (h *messagesHandler) ListChatMessages(c *echo.Context) error {
	pageParam := c.QueryParam("page")
	page, err := strconv.Atoi(pageParam)
	if err != nil || page < 1 {
		return echo.NewHTTPError(http.StatusBadRequest, "incorrect page number parameter")
	}

	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return c.NoContent(http.StatusUnauthorized)
	}

	chatID, ok := c.Get("chat_id").(uuid.UUID)
	if !ok {
		return c.NoContent(http.StatusNotFound)
	}

	dbMessages, err := h.messagesRepo.ListChatMessages(c.Request().Context(), chatID, int32(page))
	if err != nil {
		h.logger.Error("Could not load messages for chat from DB", zap.String("ChatID", chatID.String()), zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "could not load chat messages")
	}

	messages := make([]dto.Message, len(dbMessages))
	for i, m := range dbMessages {
		var readAt *time.Time
		if m.ReadAt == nil {
			readAt = nil
		} else {
			t := (*m.ReadAt).UTC()
			readAt = &t
		}
		messages[i] = dto.Message{
			ID:        m.ID.String(),
			Content:   m.Content,
			FromMe:    m.UserID == userID,
			UserID:    m.UserID.String(),
			CreatedAt: m.CreatedAt.UTC(),
			ReadAt:    readAt,
		}
	}

	resp := dto.ListChatMessagesResponse{
		Messages: messages,
	}

	return c.JSON(http.StatusCreated, resp)
}
