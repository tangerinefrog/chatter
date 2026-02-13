package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/tangerinefrog/chatter/internal/chats"
	"github.com/tangerinefrog/chatter/internal/http/dto"
	"github.com/tangerinefrog/chatter/internal/users"
	"go.uber.org/zap"
)

type chatsHandler struct {
	usersRepo *users.UsersRepository
	chatsRepo *chats.ChatsRepository
	logger    *zap.Logger
}

func NewChatsHandler(usersRepo *users.UsersRepository, chatsRepo *chats.ChatsRepository, logger *zap.Logger) *chatsHandler {
	return &chatsHandler{
		usersRepo: usersRepo,
		chatsRepo: chatsRepo,
		logger:    logger,
	}
}

func (h *chatsHandler) CreateChat(c *echo.Context) error {
	var req dto.CreateChatRequest

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

	userIDs := make([]int32, len(req.ParticipantUsernames)+1)
	// set the first participant as the current user
	userIDs[0] = userID

	for i, username := range req.ParticipantUsernames {
		u, err := h.usersRepo.GetByUsername(c.Request().Context(), username)
		if err != nil || u == nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("user with username '%s' not found", username))
		}

		if u.ID == userID {
			return echo.NewHTTPError(http.StatusBadRequest, "you cannot create a chat with yourself")
		}

		userIDs[i+1] = u.ID
	}

	var chatType chats.ChatType
	if req.IsDirect {
		alreadyExists, err := h.chatsRepo.IsDirectChatExists(c.Request().Context(), userIDs[0], userIDs[1])
		if err != nil {
			h.logger.Error("Could not check for an existing chat", zap.Int32("UserID_1", userIDs[0]), zap.Int32("UserID_2", userIDs[1]), zap.Error(err))
			return echo.NewHTTPError(http.StatusBadRequest, "could not create chat")
		}
		if alreadyExists {
			return echo.NewHTTPError(http.StatusBadRequest, "chat already exists")
		}

		chatType = chats.ChatTypeDirect
	} else {
		chatType = chats.ChatTypeGroup
	}

	chatID, err := h.chatsRepo.Create(c.Request().Context(), req.Name, chatType, userID, userIDs)
	if err != nil {
		h.logger.Error("Could not create chat between users", zap.Int32("UserID", userID), zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "could not create chat")
	}

	return c.JSON(http.StatusCreated, dto.NewChatResponse{ID: chatID})
}
