package handlers

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/tangerinefrog/chatter/internal/auth"
	"github.com/tangerinefrog/chatter/internal/http/dto"
	"github.com/tangerinefrog/chatter/internal/users"
	"go.uber.org/zap"
)

type userHandler struct {
	usersRepo *users.UsersRepository
	logger    *zap.Logger
}

func NewUserHandler(usersRepo *users.UsersRepository, logger *zap.Logger) *userHandler {
	return &userHandler{
		usersRepo: usersRepo,
		logger:    logger,
	}
}

func (h *userHandler) SignUp(c *echo.Context) error {
	var req dto.SignUpDTO

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

	hasher := auth.NewHasher()
	passwordHash, err := hasher.Hash(req.Password)
	if err != nil {
		h.logger.Error("Hashing failed", zap.Error(err))
		return c.NoContent(http.StatusInternalServerError)
	}
	_, err = h.usersRepo.Create(c.Request().Context(), req.Username, passwordHash)
	if err != nil {
		h.logger.Error("Saving new user to DB failed", zap.Error(err))
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusCreated)
}
