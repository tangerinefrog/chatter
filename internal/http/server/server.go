package server

import (
	"context"
	"net/http"
	"os"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/tangerinefrog/chatter/internal/auth/jwt"
	"github.com/tangerinefrog/chatter/internal/chats"
	"github.com/tangerinefrog/chatter/internal/http/handlers"
	mw "github.com/tangerinefrog/chatter/internal/http/middleware"
	"github.com/tangerinefrog/chatter/internal/http/validator"
	"github.com/tangerinefrog/chatter/internal/http/websockets"
	"github.com/tangerinefrog/chatter/internal/messages"
	"github.com/tangerinefrog/chatter/internal/users"
	"go.uber.org/zap"
)

type Server struct {
	addr         string
	echo         *echo.Echo
	logger       *zap.Logger
	usersRepo    *users.UsersRepository
	chatsRepo    *chats.ChatsRepository
	messagesRepo *messages.MessagesRepository
	jwtManager   *jwt.JwtManager
	hub          *websockets.Hub
}

func NewServer(
	addr string,
	logger *zap.Logger,
	usersRepo *users.UsersRepository,
	chatsRepo *chats.ChatsRepository,
	messagesRepo *messages.MessagesRepository,
	jwtManager *jwt.JwtManager,
	hub *websockets.Hub,
) *Server {
	e := echo.New()

	e.Validator = validator.New()

	e.Use(middleware.Recover())
	e.Use(middleware.RequestLogger())

	webAddr := os.Getenv("WEB_ADDR")
	if webAddr == "" {
		logger.Fatal("Web address is not defined in the .env file")
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{webAddr},
		AllowCredentials: true,
	}))

	s := Server{
		addr:         addr,
		echo:         e,
		logger:       logger,
		usersRepo:    usersRepo,
		chatsRepo:    chatsRepo,
		messagesRepo: messagesRepo,
		jwtManager:   jwtManager,
		hub:          hub,
	}

	s.registerRoutes()

	return &s
}

func (s *Server) Start(ctx context.Context) error {
	sc := echo.StartConfig{
		Address: s.addr,
	}

	err := sc.Start(ctx, s.echo)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) registerRoutes() {
	api := s.echo.Group("/api")

	authHandler := handlers.NewUserHandler(s.usersRepo, s.logger, s.jwtManager)
	auth := api.Group("/auth")
	auth.POST("/signup", authHandler.SignUp)
	auth.POST("/login", authHandler.Login)

	chatsHandler := handlers.NewChatsHandler(s.usersRepo, s.chatsRepo, s.logger)
	chats := api.Group("/chats")
	chats.Use(mw.Auth(s.jwtManager))
	chats.POST("", chatsHandler.CreateChat)
	chats.GET("", chatsHandler.ListUserChats)

	messagesHandler := handlers.NewMessagesHandler(s.messagesRepo, s.logger)
	messages := chats.Group("/:chatID/messages")
	messages.Use(mw.Chat())
	messages.GET("", messagesHandler.ListChatMessages)

	websocketsHandler := handlers.NewWebsocketsHandler(s.hub, s.logger)
	websocket := s.echo.Group("/ws")
	websocket.Use(mw.Auth(s.jwtManager))
	websocket.GET("", websocketsHandler.ServeWS)

	s.echo.GET("/health", func(c *echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
}
