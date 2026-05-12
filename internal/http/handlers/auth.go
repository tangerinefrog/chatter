package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/tangerinefrog/chatter/internal/auth/hashing"
	"github.com/tangerinefrog/chatter/internal/auth/jwt"
	"github.com/tangerinefrog/chatter/internal/http/dto"
	"github.com/tangerinefrog/chatter/internal/users"
	"go.uber.org/zap"
)

type userHandler struct {
	usersRepo  *users.UsersRepository
	logger     *zap.Logger
	jwtManager *jwt.JwtManager
}

func NewUserHandler(usersRepo *users.UsersRepository, logger *zap.Logger, jwtManager *jwt.JwtManager) *userHandler {
	return &userHandler{
		usersRepo:  usersRepo,
		logger:     logger,
		jwtManager: jwtManager,
	}
}

func (h *userHandler) SignUp(c *echo.Context) error {
	var req dto.SignUpRequest

	err := c.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	err = c.Validate(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	u, err := h.usersRepo.GetByUsername(c.Request().Context(), req.Username)
	if err != nil {
		h.logger.Error("Getting username from DB failed", zap.Error(err))
		return c.NoContent(http.StatusInternalServerError)
	}

	if u != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "username already taken")
	}

	hasher := hashing.NewHasher()
	passwordHash, err := hasher.Hash(req.Password)
	if err != nil {
		h.logger.Error("Hashing failed", zap.Error(err))
		return c.NoContent(http.StatusInternalServerError)
	}
	user, err := h.usersRepo.Create(c.Request().Context(), req.Username, passwordHash)
	if err != nil {
		h.logger.Error("Saving new user to DB failed", zap.Error(err))
		return c.NoContent(http.StatusInternalServerError)
	}

	token, expires, err := h.jwtManager.Generate(user.ID)
	if err != nil {
		h.logger.Error("Generating JWT for user failed", zap.Int32("user_id", user.ID), zap.Error(err))
		return c.NoContent(http.StatusInternalServerError)
	}

	setJwtCookie(c, token, expires)

	return c.NoContent(http.StatusOK)
}

func (h *userHandler) Login(c *echo.Context) error {
	var req dto.LogInRequest

	err := c.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	err = c.Validate(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	u, err := h.usersRepo.GetByUsername(c.Request().Context(), req.Username)
	if err != nil || u == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid username or password")
	}

	hasher := hashing.NewHasher()
	isValid, err := hasher.Verify(req.Password, u.PasswordHash)
	if err != nil {
		h.logger.Error("Hash verify failed", zap.Error(err))
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid username or password")
	}

	if !isValid {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid username or password")
	}

	token, expires, err := h.jwtManager.Generate(u.ID)
	if err != nil {
		h.logger.Error("Generating JWT for user failed", zap.Int32("user_id", u.ID), zap.Error(err))
		return c.NoContent(http.StatusInternalServerError)
	}

	setJwtCookie(c, token, expires)

	return c.NoContent(http.StatusOK)
}

func setJwtCookie(c *echo.Context, token string, expires time.Time) {
	c.SetCookie(&http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  expires,
	})
}
