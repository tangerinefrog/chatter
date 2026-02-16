package handlers

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v5"
	"github.com/tangerinefrog/chatter/internal/http/websockets"
	"go.uber.org/zap"
)

type websocketsHandler struct {
	hub      *websockets.Hub
	upgrader *websocket.Upgrader
	logger   *zap.Logger
}

func NewWebsocketsHandler(hub *websockets.Hub, logger *zap.Logger) *websocketsHandler {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	return &websocketsHandler{
		hub:      hub,
		upgrader: &upgrader,
		logger:   logger,
	}
}

func (h *websocketsHandler) ServeWS(c *echo.Context) error {
	conn, err := h.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "could not upgrade connection")
	}

	userID, ok := c.Get("user_id").(int32)
	if !ok {
		return c.NoContent(http.StatusUnauthorized)
	}

	websockets.ConnectClient(userID, conn, h.hub, h.logger)

	return nil
}
