-- 0001_init.sql
-- Initialize users table if it doesn't exist

DO $$ 
BEGIN
    -- Check if the table already exists
    IF NOT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'users') THEN
        CREATE TABLE users (
            telegram_user_id BIGINT PRIMARY KEY,
            profile          JSONB NOT NULL DEFAULT '{}'::jsonb,
            created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
            updated_at       TIMESTAMPTZ NOT NULL DEFAULT now()
        );
        
        -- Create index on telegram_user_id for faster lookups
        CREATE INDEX IF NOT EXISTS idx_users_telegram_user_id ON users(telegram_user_id);
        
        RAISE NOTICE 'Users table created successfully';
    ELSE
        RAISE NOTICE 'Users table already exists, skipping creation';
    END IF;
END $$;
