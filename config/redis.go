package config

import (
	"fmt"
	"log/slog"
	"os"
)

type RedisCred struct {
	host     string
	username string
	port     string
	password string
}

func parseRedisCred() (RedisCred, error) {
	redis_host := os.Getenv("REDIS_HOST")
	redis_password := os.Getenv("REDIS_PASSWORD")
	redis_port := os.Getenv("REDIS_PORT")
	redis_username := os.Getenv("REDIS_USERNAME")

	// checks for the .env creds;

	if redis_host == "" || redis_password == "" || redis_port == "" || redis_username == "" {
		slog.Error("Missing values for REDIS CREDENTIALS")
		return RedisCred{}, fmt.Errorf("Missing Values")
	}
	return RedisCred{
		host:     redis_host,
		username: redis_username,
		password: redis_password,
		port:     redis_port,
	}, nil
}
func CreateRedisClient() error {
	redisCredentials, err := parseRedisCred()
	if err != nil {
		return fmt.Errorf("Failed : Error creating an Redis server")
	}
	slog.Info("rdis configuration", "config", redisCredentials)
	return nil
}

