package main

import (
	"context"
	"log"
	"os"

	server "github.com/divyanshu-parihar/goFlux/cmd/api"
	config "github.com/divyanshu-parihar/goFlux/config"
	"log/slog"

	"github.com/joho/godotenv"
)

func main() {

	context.Background()

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Unable to find the .env")
	}

	slog.Info("redis configuration", "host", os.Getenv("REDIS_HOST"))
	err = config.CreateRedisClient()
	if err != nil {
		slog.Error("Redis Server Failed")
		return
	}

	slog.Info("Starting the Worker Server")
	server.CreateServer("localhost", "6379")
}
