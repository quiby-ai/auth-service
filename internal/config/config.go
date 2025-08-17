package config

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddr      string
	ShutdownTimeout time.Duration

	PGDSN string

	JWTIssuer    string
	JWTAudience  string
	JWTAccessTTL time.Duration
	JWTSecret    []byte

	TelegramBotToken string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found or error loading it: %v", err)
	}

	config := &Config{
		ServerAddr:       getEnvWithDefault("SERVER_ADDR", ":8081"),
		ShutdownTimeout:  getDurationWithDefault("SHUTDOWN_TIMEOUT_SECONDS", 30*time.Second),
		PGDSN:            os.Getenv("PG_DSN"),
		JWTIssuer:        getEnvWithDefault("JWT_ISSUER", "auth.quiby.ai"),
		JWTAudience:      getEnvWithDefault("JWT_AUDIENCE", "api.quiby.ai"),
		JWTAccessTTL:     getDurationWithDefault("JWT_ACCESS_TTL", 15*time.Minute),
		TelegramBotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
	}

	if jwtSecretB64 := os.Getenv("JWT_SECRET_B64"); jwtSecretB64 != "" {
		secret, err := base64.StdEncoding.DecodeString(jwtSecretB64)
		if err != nil {
			return nil, fmt.Errorf("invalid JWT_SECRET_B64: %w", err)
		}
		config.JWTSecret = secret
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) validate() error {
	if c.PGDSN == "" {
		return fmt.Errorf("PG_DSN environment variable is required")
	}

	if len(c.JWTSecret) == 0 {
		return fmt.Errorf("JWT_SECRET_B64 environment variable is required")
	}

	if c.TelegramBotToken == "" {
		return fmt.Errorf("TELEGRAM_BOT_TOKEN environment variable is required")
	}

	return nil
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDurationWithDefault(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if seconds, err := strconv.Atoi(value); err == nil && seconds > 0 {
			return time.Duration(seconds) * time.Second
		}
	}
	return defaultValue
}
