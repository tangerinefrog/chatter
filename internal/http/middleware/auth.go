package middleware

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/tangerinefrog/chatter/internal/auth/jwt"
)

func Auth(jwtManager *jwt.JwtManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			cookie, err := c.Request().Cookie("auth_token")
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid auth token")
			}

			claims, err := jwtManager.GetClaims(cookie.Value)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid auth token")
			}

			c.Set("user_id", claims.UserId)

			return next(c)
		}
	}
}
