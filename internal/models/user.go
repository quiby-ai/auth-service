package models

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	TelegramUserID int64           `json:"telegram_user_id"`
	Profile        json.RawMessage `json:"profile"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) UpsertUser(ctx context.Context, telegramUserID int64, profile json.RawMessage) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO users (telegram_user_id, profile)
		VALUES ($1, $2::jsonb)
		ON CONFLICT (telegram_user_id)
		DO UPDATE SET profile = EXCLUDED.profile, updated_at = now()`,
		telegramUserID, string(profile),
	)
	return err
}

func (r *UserRepository) GetUserByTelegramID(ctx context.Context, telegramUserID int64) (*User, error) {
	var user User
	err := r.db.QueryRow(ctx, `
		SELECT telegram_user_id, profile, created_at, updated_at
		FROM users WHERE telegram_user_id = $1`,
		telegramUserID,
	).Scan(&user.TelegramUserID, &user.Profile, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserProfile(ctx context.Context, telegramUserID int64) (json.RawMessage, error) {
	var profile json.RawMessage
	err := r.db.QueryRow(ctx, `
		SELECT profile FROM users WHERE telegram_user_id = $1`,
		telegramUserID,
	).Scan(&profile)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return profile, nil
}
