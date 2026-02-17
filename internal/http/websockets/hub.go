package websockets

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/tangerinefrog/chatter/internal/chats"
	"github.com/tangerinefrog/chatter/internal/messages"
	"go.uber.org/zap"
)

type Hub struct {
	mu           sync.RWMutex
	clients      map[int32]*Client
	register     chan *Client
	unregister   chan *Client
	events       chan Event
	chatsRepo    *chats.ChatsRepository
	messagesRepo *messages.MessagesRepository
	logger       *zap.Logger
}

func NewHub(chatsRepo *chats.ChatsRepository, messagesRepo *messages.MessagesRepository, logger *zap.Logger) *Hub {
	return &Hub{
		clients:      make(map[int32]*Client),
		register:     make(chan *Client),
		unregister:   make(chan *Client),
		events:       make(chan Event),
		chatsRepo:    chatsRepo,
		messagesRepo: messagesRepo,
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
		h.logger.Info("Closing WS client", zap.Int32("UserID", client.userID))
		client.close()
	}
}

func (h *Hub) handleEvent(e Event) {
	switch e.Type {
	case EventSendMessage:
		h.handleNewMessage(e.ChatID, e.SenderID, e.Content)
	}
}

func (h *Hub) handleNewMessage(chatID int32, senderID int32, message string) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	id, err := h.messagesRepo.CreateMessage(ctx, senderID, chatID, message)
	if err != nil {
		h.logger.Error("Could not save message to DB", zap.Int32("ChatID", chatID), zap.Int32("UserID", senderID), zap.Error(err))
		return
	}

	chatUsers, err := h.chatsRepo.ListUsersForChat(ctx, chatID)
	if err != nil {
		h.logger.Error("Could not get chat participants from DB", zap.Int32("ChatID", chatID), zap.Error(err))
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, u := range chatUsers {
		client, ok := h.clients[u.ID]
		if ok {
			msg := Event{
				Type:      EventNewMessage,
				Content:   message,
				SenderID:  senderID,
				ChatID:    chatID,
				MessageID: id,
				FromMe:    u.ID == senderID,
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
		h.logger.Debug("Could not send message to client, skip", zap.Int32("UserID", c.userID))
	}
}
