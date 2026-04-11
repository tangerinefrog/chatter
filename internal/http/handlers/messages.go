package handlers

import (
	"net/http"
	"sort"
	"strconv"
	"time"

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

	userID, ok := c.Get("user_id").(int32)
	if !ok {
		return c.NoContent(http.StatusUnauthorized)
	}

	chatID, ok := c.Get("chat_id").(int32)
	if !ok {
		return c.NoContent(http.StatusNotFound)
	}

	dbMessages, err := h.messagesRepo.ListChatMessages(c.Request().Context(), chatID, int32(page))
	if err != nil {
		h.logger.Error("Could not load messages for chat from DB", zap.Int32("ChatID", chatID), zap.Error(err))
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
			ID:        m.ID,
			Content:   m.Content,
			FromMe:    m.UserID == userID,
			UserID:    m.UserID,
			CreatedAt: m.CreatedAt.UTC(),
			ReadAt:    readAt,
		}
	}

	sort.Slice(messages, func(i, j int) bool {
		return messages[i].ID < messages[j].ID
	})

	resp := dto.ListChatMessagesResponse{
		Messages: messages,
	}

	return c.JSON(http.StatusCreated, resp)
}

func (h *messagesHandler) MarkMessagesAsRead(c *echo.Context) error {
	var req dto.ReadMessageRequest

	err := c.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	err = c.Validate(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	userID, ok := c.Get("user_id").(int32)
	if !ok {
		return c.NoContent(http.StatusUnauthorized)
	}

	chatID, ok := c.Get("chat_id").(int32)
	if !ok {
		return c.NoContent(http.StatusNotFound)
	}

	err = h.messagesRepo.MarkMessagesAsRead(c.Request().Context(), req.MessageID, chatID, userID)
	if err != nil {
		h.logger.Error("Could not mark messages as read for chat in DB", zap.Int32("ChatID", chatID), zap.Int64("MessageID", req.MessageID), zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "could not mark messages as read")
	}

	return c.NoContent(http.StatusOK)
}
