package config

import (
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
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
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("/")

	viper.AutomaticEnv()

	viper.BindEnv("server.addr", "SERVER_ADDR")
	viper.BindEnv("server.shutdown_timeout_seconds", "SHUTDOWN_TIMEOUT_SECONDS")
	viper.BindEnv("jwt.issuer", "JWT_ISSUER")
	viper.BindEnv("jwt.audience", "JWT_AUDIENCE")
	viper.BindEnv("jwt.access_ttl_minutes", "JWT_ACCESS_TTL")

	viper.BindEnv("PG_DSN")
	viper.BindEnv("JWT_SECRET")
	viper.BindEnv("TELEGRAM_BOT_TOKEN")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("No config.toml file found, using environment variables and defaults")
		} else {
			log.Printf("Error reading config file: %v", err)
		}
	}

	config := &Config{
		ServerAddr:       viper.GetString("server.addr"),
		ShutdownTimeout:  time.Duration(viper.GetInt("server.shutdown_timeout_seconds")) * time.Second,
		PGDSN:            viper.GetString("PG_DSN"),
		JWTIssuer:        viper.GetString("jwt.issuer"),
		JWTAudience:      viper.GetString("jwt.audience"),
		JWTAccessTTL:     time.Duration(viper.GetInt("jwt.access_ttl_minutes")) * time.Minute,
		TelegramBotToken: viper.GetString("TELEGRAM_BOT_TOKEN"),
	}

	if jwtSecret := viper.GetString("JWT_SECRET"); jwtSecret != "" {
		secret, err := base64.StdEncoding.DecodeString(jwtSecret)
		if err != nil {
			return nil, fmt.Errorf("invalid JWT_SECRET: %w", err)
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
		return fmt.Errorf("JWT_SECRET environment variable is required")
	}

	if c.TelegramBotToken == "" {
		return fmt.Errorf("TELEGRAM_BOT_TOKEN environment variable is required")
	}

	return nil
}

func GetString(key string) string {
	return viper.GetString(key)
}

func GetInt(key string) int {
	return viper.GetInt(key)
}

func GetDuration(key string) time.Duration {
	return viper.GetDuration(key)
}

func IsSet(key string) bool {
	return viper.IsSet(key)
}
