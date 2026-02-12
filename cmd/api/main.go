package main

import (
	"context"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/tangerinefrog/chatter/internal/auth/jwt"
	"github.com/tangerinefrog/chatter/internal/db"
	"github.com/tangerinefrog/chatter/internal/http/server"
	"github.com/tangerinefrog/chatter/internal/users"
	"go.uber.org/zap"
)

func main() {
	logger := zap.Must(zap.NewProduction())
	defer logger.Sync()

	err := godotenv.Load()
	if err != nil {
		logger.Fatal("could not load .env file", zap.Error(err))
	}

	connStr := os.Getenv("DB_CONN")
	if connStr == "" {
		logger.Fatal("DB connection string is not defined in the .env file")
	}

	addr := os.Getenv("SRV_ADDR")
	if addr == "" {
		logger.Fatal("Server address is not defined in the .env file")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		logger.Fatal("JWT secret is not defined in the .env file")
	}

	jwtManager := jwt.NewJwtManager(jwtSecret, 1*time.Hour)

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		logger.Fatal("Error creating a DB pool")
	}
	defer pool.Close()

	queries := db.New(pool)

	usersRepo := users.NewRepository(queries)

	srv := server.NewServer(addr, logger, usersRepo, jwtManager)
	srv.Start()
}
