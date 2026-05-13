package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

func Chat() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			param := c.Param("chatID")

			id, err := uuid.Parse(param)
			if err != nil {
				return echo.NewHTTPError(http.StatusNotFound, "chat not found")
			}
			c.Set("chat_id", id)

			return next(c)
		}
	}
}
