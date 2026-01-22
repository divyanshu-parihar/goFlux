package config

import (
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

// TestParseRedisCred tests the credential parsing logic.
// It does NOT require a running Redis.
func TestParseRedisCred(t *testing.T) {
	_ = godotenv.Load("../.env")

	// Save current env to restore later
	oldHost := os.Getenv("REDIS_HOST")
	oldPort := os.Getenv("REDIS_PORT")
	oldPass := os.Getenv("REDIS_PASSWORD")
	oldUser := os.Getenv("REDIS_USERNAME")
	defer func() {
		os.Setenv("REDIS_HOST", oldHost)
		os.Setenv("REDIS_PORT", oldPort)
		os.Setenv("REDIS_PASSWORD", oldPass)
		os.Setenv("REDIS_USERNAME", oldUser)
	}()

	t.Run("Success", func(t *testing.T) {
		os.Setenv("REDIS_HOST", os.Getenv("TEST_REDIS_HOST"))
		os.Setenv("REDIS_PORT", os.Getenv("TEST_REDIS_PORT"))
		os.Setenv("REDIS_PASSWORD", os.Getenv("TEST_REDIS_PASSWORD"))
		os.Setenv("REDIS_USERNAME", os.Getenv("TEST_REDIS_USERNAME"))

		cred, err := parseRedisCred()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if cred.host != os.Getenv("TEST_REDIS_HOST") {
			t.Errorf("Expected host %s, got %s", os.Getenv("TEST_REDIS_HOST"), cred.host)
		}
	})

	t.Run("MissingValues", func(t *testing.T) {
		os.Unsetenv("REDIS_HOST")
		// Assuming other vars are still set from previous run or empty
		// We just need one to be missing to trigger error

		_, err := parseRedisCred()
		if err == nil {
			t.Error("Expected error for missing REDIS_HOST, got nil")
		}
	})
}

// TestRedisOperations tests the CRUD functions.
// This REQUIRES a running Redis instance or it will be skipped.
func TestRedisOperations(t *testing.T) {
	_ = godotenv.Load("../.env")

	os.Setenv("REDIS_HOST", os.Getenv("TEST_REDIS_HOST"))
	os.Setenv("REDIS_PORT", os.Getenv("TEST_REDIS_PORT"))
	os.Setenv("REDIS_PASSWORD", os.Getenv("TEST_REDIS_PASSWORD"))
	os.Setenv("REDIS_USERNAME", os.Getenv("TEST_REDIS_USERNAME"))

	client, err := CreateRedisClient()
	if err != nil {
		t.Skipf("Skipping integration test: could not create client: %v", err)
	}

	// Ping to check if actually reachable
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		t.Skipf("Skipping integration test: Redis not reachable: %v", err)
	}
	defer client.Close()

	t.Run("CreateAndGetHValue", func(t *testing.T) {
		key := "test_queue"
		field := "job_1"
		value := "processing"

		// Test Create
		// Note: The function signature is CreateRedisHValue(ctx, client, key, value)
		// But internally it does HSet(ctx, key, value).
		// HSet expects key, then field-value pairs.
		// Wait, the code says: client.HSet(ctx, key, value).Result()
		// If value is a struct or map, HSet works. If value is just a string... HSet expects pairs.
		// Let's look at code: func CreateRedisHValue(..., value interface{})
		// If I pass "processing" as value, HSet("test_queue", "processing") -> Error? "WRONGTYPE Operation against a key holding the wrong kind of value" or "ERR wrong number of arguments for 'hset' command"
		// HSet requires key field value [field value ...].
		// The implementation seems to assume `value` contains the field(s).

		// Let's pass a map which HSet supports
		input := map[string]interface{}{
			field: value,
		}

		err := CreateRedisHValue(ctx, client, key, input)
		if err != nil {
			t.Fatalf("Failed to create value: %v", err)
		}

		// Test Get
		got, err := GetRedisHValue(ctx, client, key, field)
		if err != nil {
			t.Fatalf("Failed to get value: %v", err)
		}
		if got != value {
			t.Errorf("Expected %s, got %s", value, got)
		}
	})

	t.Run("GetAllHValue", func(t *testing.T) {
		key := "test_queue_all"
		data := map[string]interface{}{
			"f1": "v1",
			"f2": "v2",
		}

		err := CreateRedisHValue(ctx, client, key, data)
		if err != nil {
			t.Fatalf("Failed to setup data: %v", err)
		}

		result, err := GetRedisAllHValue(ctx, client, key)
		if err != nil {
			t.Fatalf("Failed to get all values: %v", err)
		}

		if len(result) != 2 {
			t.Errorf("Expected 2 items, got %d", len(result))
		}
		if result["f1"] != "v1" {
			t.Errorf("Expected v1, got %s", result["f1"])
		}
	})
}
