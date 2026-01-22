package main

import (
	server "github.com/divyanshu-parihar/goFlux/cmd/api"
	"log"
	"log/slog"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Unable to find the .env")
	}
	slog.Info("Starting the Worker Server")
	server.CreateServer("localhost", "6379")
}
