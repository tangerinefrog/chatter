package websockets

import (
	"context"
	"encoding/json"
	"time"

	"github.com/tangerinefrog/chatter/internal/chats"
	"github.com/tangerinefrog/chatter/internal/messages"
	"go.uber.org/zap"
)

type Hub struct {
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
			h.clients[client.userID] = client
		case client := <-h.unregister:
			_, ok := h.clients[client.userID]
			if ok {
				close(client.send)
				delete(h.clients, client.userID)
			}
		case e := <-h.events:
			h.handleEvent(e)
		}
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
	for _, u := range chatUsers {
		client, ok := h.clients[u.ID]
		if ok {
			msg := Event{
				Type:      EventNewMessage,
				Content:   message,
				SenderID:  senderID,
				ChatID:    chatID,
				MessageID: id,
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

	c.send <- messageBytes
}
