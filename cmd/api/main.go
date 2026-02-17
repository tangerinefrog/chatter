package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/tangerinefrog/chatter/internal/auth/jwt"
	"github.com/tangerinefrog/chatter/internal/chats"
	"github.com/tangerinefrog/chatter/internal/http/server"
	"github.com/tangerinefrog/chatter/internal/http/websockets"
	"github.com/tangerinefrog/chatter/internal/messages"
	"github.com/tangerinefrog/chatter/internal/users"
	"go.uber.org/zap"
)

type config struct {
	dbConn     string
	serverAddr string
	jwtSecret  string
	webAddr    string
}

func main() {
	logger := zap.Must(zap.NewProduction())
	defer logger.Sync()
	err := run(logger)
	if err != nil {
		logger.Fatal("Server error", zap.Error(err))
	}
}

func run(logger *zap.Logger) error {
	err := godotenv.Load()
	if err != nil {
		logger.Warn("could not load .env file, using environment variables", zap.Error(err))
	}

	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	pool, err := initDB(ctx, cfg.dbConn, logger)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer pool.Close()

	usersRepo := users.NewRepository(pool)
	chatsRepo := chats.NewRepository(pool)
	messagesRepo := messages.NewRepository(pool)

	hub := websockets.NewHub(chatsRepo, messagesRepo, logger)
	go hub.Run()

	jwtManager := jwt.NewJwtManager(cfg.jwtSecret, 1*time.Hour)
	srv := server.NewServer(cfg.serverAddr, logger, usersRepo, chatsRepo, messagesRepo, jwtManager, hub)

	err = srv.Start(ctx)
	if err != nil {
		return err
	}

	logger.Info("Shutting down server...")
	

	return nil
}

func loadConfig() (*config, error) {
	cfg := &config{
		dbConn:     os.Getenv("DB_CONN"),
		serverAddr: os.Getenv("SRV_ADDR"),
		jwtSecret:  os.Getenv("JWT_SECRET"),
		webAddr:    os.Getenv("WEB_ADDR"),
	}

	var missing []string
	if cfg.dbConn == "" {
		missing = append(missing, "DB_CONN")
	}
	if cfg.serverAddr == "" {
		missing = append(missing, "SRV_ADDR")
	}
	if cfg.jwtSecret == "" {
		missing = append(missing, "JWT_SECRET")
	}
	if cfg.webAddr == "" {
		missing = append(missing, "WEB_ADDR")
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required environment variables: %v", missing)
	}

	return cfg, nil
}

func initDB(ctx context.Context, connStr string, logger *zap.Logger) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}
