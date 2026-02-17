package websockets

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Client struct {
	userID int32
	conn   *websocket.Conn
	send   chan []byte
	hub    *Hub
	logger *zap.Logger
	mu     sync.Mutex
	once   sync.Once
}

func ConnectClient(userID int32, conn *websocket.Conn, hub *Hub, logger *zap.Logger) {
	c := &Client{
		userID: userID,
		conn:   conn,
		send:   make(chan []byte),
		hub:    hub,
		logger: logger,
	}

	hub.register <- c

	go c.readPump()
	go c.writePump()
}

func (c *Client) readPump() {
	defer c.close()

	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			c.logger.Error("Could not read WS message", zap.Int32("UserID", c.userID), zap.Error(err))
			break
		}

		var event Event
		err = json.Unmarshal(data, &event)
		if err != nil {
			c.logger.Error("Could not deserialize WS message", zap.String("Event", string(data)), zap.Int32("UserID", c.userID), zap.Error(err))
			break
		}

		c.logger.Info("Got WS message", zap.Any("Event", event))
		event.SenderID = c.userID

		c.hub.events <- event
	}
}

func (c *Client) writePump() {
	defer c.close()
	for data := range c.send {
		c.mu.Lock()
		err := c.conn.WriteMessage(websocket.TextMessage, data)
		c.mu.Unlock()
		if err != nil {
			return
		}
	}
}

func (c *Client) close() {
	c.once.Do(func() {
		c.hub.unregister <- c
		c.mu.Lock()
		defer c.mu.Unlock()
		c.conn.WriteMessage(websocket.CloseMessage, []byte{})
		c.conn.Close()
	})
}
