package middleware

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
)

func Chat() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			param := c.Param("chatID")

			chatID, err := strconv.Atoi(param)
			if err != nil {
				return echo.NewHTTPError(http.StatusNotFound, "chat not found")
			}

			c.Set("chat_id", int32(chatID))

			return next(c)
		}
	}
}
