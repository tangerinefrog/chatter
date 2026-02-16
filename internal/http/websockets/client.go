package websockets

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Client struct {
	userID int32
	conn   *websocket.Conn
	send   chan []byte
	hub    *Hub
	logger *zap.Logger
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
		err := c.conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			return
		}
	}
}

func (c *Client) close() {
	c.hub.unregister <- c
	c.conn.WriteMessage(websocket.CloseMessage, []byte{})
	c.conn.Close()
}
