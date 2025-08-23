package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresConnection(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres config: %w", err)
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL")

	if err := InitDB(ctx, pool); err != nil {
		log.Printf("Warning: Failed to initialize database schema: %v", err)
	}

	return pool, nil
}

func InitDB(ctx context.Context, pool *pgxpool.Pool) error {
	migrationSQL, err := os.ReadFile("/init.sql")
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	_, err = pool.Exec(ctx, string(migrationSQL))
	if err != nil {
		return fmt.Errorf("failed to run database migration: %w", err)
	}

	log.Println("Database schema initialized successfully")
	return nil
}

func ClosePostgresConnection(pool *pgxpool.Pool) {
	if pool != nil {
		pool.Close()
		log.Println("PostgreSQL connection closed")
	}
}
