-- Add email and password fields to users table
ALTER TABLE users 
ADD COLUMN IF NOT EXISTS email TEXT UNIQUE,
ADD COLUMN IF NOT EXISTS password_hash TEXT;

-- Make telegram_id nullable (users can register without Telegram)
ALTER TABLE users 
ALTER COLUMN telegram_id DROP NOT NULL;

-- Add index for email
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Add constraint: user must have either telegram_id or email
ALTER TABLE users 
ADD CONSTRAINT users_telegram_or_email CHECK (
    telegram_id IS NOT NULL OR email IS NOT NULL
);

