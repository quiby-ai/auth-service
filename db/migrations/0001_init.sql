-- 0001_init.sql
CREATE TABLE IF NOT EXISTS users (
  telegram_user_id BIGINT PRIMARY KEY,
  profile          JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- helpful GIN index if you ever filter by profile keys later (not required now)
-- CREATE INDEX IF NOT EXISTS users_profile_gin ON users USING GIN (profile);
