package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("could not load .env file: %w", err)
	}

	addr := os.Getenv("SRV_ADDR")
	if addr == "" {
		log.Fatalf("server address is not defined in the .env file")
	}

	e := echo.New()

	if err := e.Start(addr); err != nil {
		log.Fatalf("failed to start the server: %v", err)
	}
}

