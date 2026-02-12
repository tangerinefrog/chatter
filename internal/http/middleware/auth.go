package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/tangerinefrog/chatter/internal/auth/jwt"
)

func Auth(jwtManager *jwt.JwtManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing auth token")
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid auth header")
			}

			claims, err := jwtManager.GetClaims(parts[1])
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid auth token")
			}

			c.Set("user_id", claims.UserId)

			return next(c)
		}
	}
}
