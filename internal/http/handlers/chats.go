package handlers

import (
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/google/uuid"
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

	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return c.NoContent(http.StatusUnauthorized)
	}

	userIDs := make([]uuid.UUID, len(req.ParticipantUsernames)+1)
	// set the first participant as the current user
	userIDs[0] = userID

	for i, username := range req.ParticipantUsernames {
		u, err := h.usersRepo.GetByUsername(c.Request().Context(), username)
		if err != nil || u == nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("user with username '%s' not found", username))
		}

		if u.ID.Bytes == userID {
			return echo.NewHTTPError(http.StatusBadRequest, "you cannot create a chat with yourself")
		}

		userIDs[i+1] = u.ID.Bytes
	}

	var chatType chats.ChatType
	if req.IsDirect {
		alreadyExists, err := h.chatsRepo.IsDirectChatExists(c.Request().Context(), userIDs[0], userIDs[1])
		if err != nil {
			h.logger.Error("Could not check for an existing chat", zap.String("UserID_1", userIDs[0].String()), zap.String("UserID_2", userIDs[1].String()), zap.Error(err))
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
		h.logger.Error("Could not create chat between users", zap.String("UserID", userID.String()), zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "could not create chat")
	}

	return c.JSON(http.StatusCreated, dto.NewChatResponse{ID: chatID.String()})
}

func (h *chatsHandler) ListUserChats(c *echo.Context) error {
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return c.NoContent(http.StatusUnauthorized)
	}

	dbChats, err := h.chatsRepo.ListChatsForUser(c.Request().Context(), userID)
	if err != nil {
		h.logger.Error("Could not get chats for user", zap.String("UserID", userID.String()), zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "could not get chats")
	}

	chats := make([]dto.Chat, len(dbChats))

	for i, c := range dbChats {
		participants := make([]dto.ChatParticipant, len(c.Participants))
		chatName := c.Name
		for i, p := range c.Participants {
			participants[i] = dto.ChatParticipant{
				ID:       p.ID.String(),
				Username: p.Username,
			}
			if chatName == "" && p.ID != userID {
				chatName = p.Username
			}
		}

		lastMessageDate := &c.LastMessageDate
		epoch := time.Time{}
		if (*lastMessageDate).Equal(epoch) {
			lastMessageDate = nil
		}

		chats[i] = dto.Chat{
			ID:                  c.ID.String(),
			Type:                c.Type,
			Name:                chatName,
			Participants:        participants,
			LastMessage:         c.LastMessage,
			LastMessageDate:     lastMessageDate,
			UnreadMessagesCount: c.UnreadMessagesCount,
		}
	}

	sort.Slice(chats, func(i, j int) bool {
		if chats[i].LastMessageDate == nil {
			return false
		}
		if chats[j].LastMessageDate == nil {
			return true
		}

		return chats[i].LastMessageDate.After(*chats[j].LastMessageDate)
	})

	resp := dto.ListChatsForUserResponse{
		Chats: chats,
	}

	return c.JSON(http.StatusOK, resp)
}
