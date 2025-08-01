package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken        string
	JWTSecret       string
	Port            string
	ShutdownTimeout time.Duration
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found or error loading it: %v", err)
	}

	config := &Config{
		BotToken:        os.Getenv("BOT_TOKEN"),
		JWTSecret:       os.Getenv("JWT_SECRET"),
		Port:            getEnvWithDefault("PORT", "8081"),
		ShutdownTimeout: getDurationWithDefault("SHUTDOWN_TIMEOUT_SECONDS", 30*time.Second),
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) validate() error {
	if c.BotToken == "" {
		return fmt.Errorf("BOT_TOKEN environment variable is required")
	}

	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET environment variable is required")
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
