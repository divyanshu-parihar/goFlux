package config

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/redis/go-redis/v9"
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

// CreateRedisClient initialized the redis client required to store the queue
func CreateRedisClient() (*redis.Client, error) {
	redisCredentials, err := parseRedisCred()
	if err != nil {
		return &redis.Client{}, fmt.Errorf("Failed : Error creating an Redis server")
	}
	slog.Info("rdis configuration", "config", redisCredentials)
	return redis.NewClient(&redis.Options{
		Addr:     redisCredentials.host + redisCredentials.port,
		Password: redisCredentials.password, // no password
		DB:       0,                         // use default DB
		Protocol: 2,
	}), nil
}

// functions for the redis queue ops

// CRUD

func CreateRedisHValue(ctx context.Context, client *redis.Client, key string, value interface{}) error {
	_, err := client.HSet(ctx, key, value).Result()

	if err != nil {
		slog.Error("Error inserting into Redis : ", key, value)
		return err
	}
	return nil
}

func GetRedisHValue(ctx context.Context, client *redis.Client, hsetKey, value string) (string, error) {
	result, err := client.HGet(ctx, hsetKey, value).Result()

	if err != nil {
		slog.Error("Error getting into Redis : ", hsetKey, value)
		return "", err
	}
	return result, nil
}

func GetRedisAllHValue(ctx context.Context, client *redis.Client, hsetKey, value string) (map[string]string, error) {
	result, err := client.HGetAll(ctx, hsetKey).Result()

	if err != nil {
		slog.Error("Error getting into Redis : ", hsetKey, value)
		return make(map[string]string), err
	}
	return result, nil
}
