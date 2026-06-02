package websockets

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/tangerinefrog/chatter/internal/chats"
	"github.com/tangerinefrog/chatter/internal/files"
	"github.com/tangerinefrog/chatter/internal/messages"
	"go.uber.org/zap"
)

type Hub struct {
	mu           sync.RWMutex
	clients      map[uuid.UUID]*Client
	register     chan *Client
	unregister   chan *Client
	events       chan Event
	chatsRepo    *chats.ChatsRepository
	messagesRepo *messages.MessagesRepository
	filesRepo    *files.FilesRepository
	logger       *zap.Logger
}

func NewHub(chatsRepo *chats.ChatsRepository, messagesRepo *messages.MessagesRepository, filesRepo *files.FilesRepository, logger *zap.Logger) *Hub {
	return &Hub{
		clients:      make(map[uuid.UUID]*Client),
		register:     make(chan *Client),
		unregister:   make(chan *Client),
		events:       make(chan Event),
		chatsRepo:    chatsRepo,
		messagesRepo: messagesRepo,
		filesRepo:    filesRepo,
		logger:       logger,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.userID] = client
			h.mu.Unlock()
		case client := <-h.unregister:
			h.mu.Lock()
			_, ok := h.clients[client.userID]
			if ok {
				close(client.send)
				delete(h.clients, client.userID)
			}
			h.mu.Unlock()
		case e := <-h.events:
			h.handleEvent(e)
		}
	}
}

func (h *Hub) Shutdown() {
	clients := make([]*Client, 0)
	h.mu.Lock()
	for _, client := range h.clients {
		clients = append(clients, client)
	}
	h.mu.Unlock()
	h.logger.Info("WS Clients", zap.Any("Clients", clients))

	for _, client := range clients {
		h.logger.Info("Closing WS client", zap.String("UserID", client.userID.String()))
		client.close()
	}
}

func (h *Hub) handleEvent(e Event) {
	switch e.Type {
	case EventSendMessage:
		h.handleNewMessage(e.ChatID, e.SenderID, e.Content, e.FileIDs)
	case EventReadMessage:
		h.handleReadMessage(e.ChatID, e.SenderID, e.MessageID)
	}
}

func (h *Hub) handleNewMessage(chatID uuid.UUID, senderID uuid.UUID, message string, fileIDs []uuid.UUID) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	id, err := h.messagesRepo.CreateMessage(ctx, senderID, chatID, message)
	if err != nil {
		h.logger.Error("Could not save message to DB", zap.String("ChatID", chatID.String()), zap.String("UserID", senderID.String()), zap.Error(err))
		return
	}

	var fileInfos []FileInfo
	for _, fileID := range fileIDs {
		if fileID != uuid.Nil {
			err := h.filesRepo.LinkFileToMessage(ctx, fileID, id)
			if err != nil {
				h.logger.Error("Could not link file to message", zap.String("FileID", fileID.String()), zap.String("MessageID", id.String()), zap.Error(err))
				continue
			}

			file, err := h.filesRepo.GetFileByID(ctx, fileID)
			if err != nil {
				h.logger.Error("Could not get file details", zap.String("FileID", fileID.String()), zap.Error(err))
				continue
			}

			fileInfos = append(fileInfos, FileInfo{
				ID:        file.ID.String(),
				Name:      file.FileName,
				MimeType:  file.MimeType,
				SizeBytes: file.SizeBytes,
			})
		}
	}

	chatUsers, err := h.chatsRepo.ListUsersForChat(ctx, chatID)
	if err != nil {
		h.logger.Error("Could not get chat participants from DB", zap.String("ChatID", chatID.String()), zap.Error(err))
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, u := range chatUsers {
		client, ok := h.clients[u.ID.Bytes]
		if ok {
			msg := Event{
				Type:      EventNewMessage,
				Content:   message,
				SenderID:  senderID,
				ChatID:    chatID,
				MessageID: id,
				Files:     fileInfos,
				FromMe:    u.ID.Bytes == senderID,
			}
			h.sendMessageToClient(msg, client)
		}
	}
}

func (h *Hub) sendMessageToClient(m Event, c *Client) {
	messageBytes, err := json.Marshal(m)
	if err != nil {
		h.logger.Error("Could not serialize WS message", zap.Any("Message", m))
		return
	}

	select {
	case c.send <- messageBytes:
	default:
		h.logger.Debug("Could not send message to client, skip", zap.String("UserID", c.userID.String()))
	}
}

func (h *Hub) handleReadMessage(chatID uuid.UUID, senderID uuid.UUID, messageID uuid.UUID) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	if err := h.messagesRepo.MarkMessagesAsRead(ctx, messageID, chatID, senderID); err != nil {
		h.logger.Error("Could not mark messages as read in DB", zap.String("ChatID", chatID.String()), zap.String("MessageID", messageID.String()), zap.String("UserID", senderID.String()), zap.Error(err))
	}

	chatUsers, err := h.chatsRepo.ListUsersForChat(ctx, chatID)
	if err != nil {
		h.logger.Error("Could not get chat participants from DB", zap.String("ChatID", chatID.String()), zap.Error(err))
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, u := range chatUsers {
		if u.ID.Bytes == senderID {
			continue
		}

		client, ok := h.clients[u.ID.Bytes]
		if ok {
			msg := Event{
				Type:      EventReadMessage,
				SenderID:  senderID,
				ChatID:    chatID,
				MessageID: messageID,
				ReadAt:    time.Now().UTC(),
			}
			h.sendMessageToClient(msg, client)
		}
	}
}
