-- Add SSO specific columns and metadata for internal (Pegawai) and external (Mitra)
ALTER TABLE users ADD COLUMN IF NOT EXISTS user_type VARCHAR(20) NOT NULL DEFAULT 'internal';
ALTER TABLE users ADD COLUMN IF NOT EXISTS nik VARCHAR(20);
ALTER TABLE users ADD COLUMN IF NOT EXISTS metadata JSONB DEFAULT '{}'::jsonb;

CREATE INDEX IF NOT EXISTS idx_users_user_type ON users (user_type);
CREATE INDEX IF NOT EXISTS idx_users_nik ON users (nik);
