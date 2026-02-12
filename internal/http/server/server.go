package server

import (
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/tangerinefrog/chatter/internal/auth/jwt"
	"github.com/tangerinefrog/chatter/internal/http/handlers"
	"github.com/tangerinefrog/chatter/internal/http/validator"
	"github.com/tangerinefrog/chatter/internal/users"
	"go.uber.org/zap"
)

type Server struct {
	addr       string
	echo       *echo.Echo
	logger     *zap.Logger
	usersRepo  *users.UsersRepository
	jwtManager *jwt.JwtManager
}

func NewServer(addr string, logger *zap.Logger, usersRepo *users.UsersRepository) *Server {
	e := echo.New()

	e.Validator = validator.New()

	e.Use(middleware.Recover())
	e.Use(middleware.RequestLogger())
	e.Use(middleware.CORS("*"))

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		logger.Fatal("JWT secret is not defined in the .env file")
	}

	jwtManager := jwt.NewJwtManager(jwtSecret, 1*time.Hour)

	s := Server{
		addr:       addr,
		echo:       e,
		logger:     logger,
		usersRepo:  usersRepo,
		jwtManager: jwtManager,
	}

	s.registerRoutes()

	return &s
}

func (s *Server) Start() {
	err := s.echo.Start(s.addr)
	if err != nil {
		s.logger.Fatal("Error during server start", zap.Error(err))
	}
}

func (s *Server) registerRoutes() {
	api := s.echo.Group("/api")

	authHandler := handlers.NewUserHandler(s.usersRepo, s.logger, s.jwtManager)
	auth := api.Group("/auth")
	auth.POST("/signup", authHandler.SignUp)
	auth.GET("/login", authHandler.Login)

	s.echo.GET("/health", func(c *echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
}
