package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	"github.com/tangerinefrog/chatter/internal/auth/jwt"
	"github.com/tangerinefrog/chatter/internal/chats"
	"github.com/tangerinefrog/chatter/internal/crypto"
	"github.com/tangerinefrog/chatter/internal/files"
	"github.com/tangerinefrog/chatter/internal/http/server"
	"github.com/tangerinefrog/chatter/internal/http/websockets"
	"github.com/tangerinefrog/chatter/internal/messages"
	"github.com/tangerinefrog/chatter/internal/storage"
	"github.com/tangerinefrog/chatter/internal/users"
	"github.com/tangerinefrog/chatter/sql"
	"go.uber.org/zap"
)

type config struct {
	dbConn        string
	serverAddr    string
	jwtSecret     string
	encryptionKey []byte
	s3Endpoint    string
	s3Bucket      string
	s3Region      string
	s3Username    string
	s3Password    string
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

	pool, err := initDB(ctx, cfg.dbConn)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer pool.Close()

	cipher, err := crypto.NewCipher(cfg.encryptionKey)

	usersRepo := users.NewRepository(pool)
	chatsRepo := chats.NewRepository(pool, cipher)
	messagesRepo := messages.NewRepository(pool, cipher)
	filesRepo := files.NewRepository(pool)

	s3Storage, err := storage.NewS3Storage(cfg.s3Endpoint, cfg.s3Bucket, cfg.s3Region, cfg.s3Username, cfg.s3Password)
	if err != nil {
		return fmt.Errorf("failed to initialize S3 storage: %w", err)
	}

	filesService := files.NewFileService(filesRepo, s3Storage, cipher)

	hub := websockets.NewHub(chatsRepo, messagesRepo, logger)
	go hub.Run()

	jwtManager := jwt.NewJwtManager(cfg.jwtSecret, 1*time.Hour)
	srv := server.NewServer(
		cfg.serverAddr,
		logger,
		usersRepo,
		chatsRepo,
		messagesRepo,
		jwtManager,
		hub,
		filesService,
	)

	err = srv.Start(ctx)
	if err != nil {
		return err
	}

	logger.Info("Shutting down server...")

	return nil
}

func loadConfig() (*config, error) {
	encryptionKeyStr := os.Getenv("ENCRYPTION_KEY")
	var encryptionKey []byte
	var err error

	if encryptionKeyStr != "" {
		encryptionKey, err = hex.DecodeString(encryptionKeyStr)
		if err != nil {
			return nil, fmt.Errorf("invalid ENCRYPTION_KEY format (must be hex-encoded): %w", err)
		}
		if len(encryptionKey) != 32 {
			return nil, fmt.Errorf("ENCRYPTION_KEY must be 32 bytes when decoded, got %d", len(encryptionKey))
		}
	}

	cfg := &config{
		dbConn:        os.Getenv("DB_CONN"),
		serverAddr:    os.Getenv("SRV_ADDR"),
		jwtSecret:     os.Getenv("JWT_SECRET"),
		s3Endpoint:    os.Getenv("S3_ENDPOINT"),
		s3Bucket:      os.Getenv("S3_BUCKET"),
		s3Region:      os.Getenv("S3_REGION"),
		s3Username:    os.Getenv("S3_USERNAME"),
		s3Password:    os.Getenv("S3_PASSWORD"),
		encryptionKey: encryptionKey,
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
	if len(cfg.encryptionKey) == 0 {
		missing = append(missing, "ENCRYPTION_KEY")
	}
	if cfg.s3Endpoint == "" {
		missing = append(missing, "S3_ENDPOINT")
	}
	if cfg.s3Bucket == "" {
		missing = append(missing, "S3_BUCKET")
	}
	if cfg.s3Region == "" {
		missing = append(missing, "S3_REGION")
	}
	if cfg.s3Username == "" {
		missing = append(missing, "S3_USERNAME")
	}
	if cfg.s3Password == "" {
		missing = append(missing, "S3_PASSWORD")
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required environment variables: %v", missing)
	}

	return cfg, nil
}

func initDB(ctx context.Context, connStr string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	goose.SetBaseFS(sql.EmbedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to set database dialect: %w", err)
	}

	db := stdlib.OpenDBFromPool(pool)

	if err := goose.UpContext(ctx, db, "migrations"); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to apply database migrations: %w", err)
	}

	return pool, nil
}
