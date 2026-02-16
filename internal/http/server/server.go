package server

import (
	"net/http"

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
	e.Use(middleware.CORS("*"))

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
	messages.POST("", messagesHandler.CreateMessage)

	websocketsHandler := handlers.NewWebsocketsHandler(s.hub, s.logger)
	websocket := s.echo.Group("/ws")
	websocket.Use(mw.Auth(s.jwtManager))
	websocket.GET("", websocketsHandler.ServeWS)

	s.echo.GET("/health", func(c *echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
}
