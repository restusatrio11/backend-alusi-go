-- Add manual authentication fields to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS username VARCHAR(100);
ALTER TABLE users ADD COLUMN IF NOT EXISTS password_hash VARCHAR(255);
ALTER TABLE users ALTER COLUMN sso_sub DROP NOT NULL;

-- Create unique index on username
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username ON users (username) WHERE username IS NOT NULL;

-- Seed default initial admin account if not already present
-- Default credentials: username = 'admin', password = 'AdminBPS1200!' (bcrypt hash)
INSERT INTO users (sso_sub, username, nip, nama, email, status, password_hash, user_type)
VALUES (
    'manual:admin',
    'admin',
    '199501012020011001',
    'Administrator ALUSI',
    'admin@bps.go.id',
    'active',
    '$2a$10$CtL0O5Pagh8WXov1KJIqVOmf7w0WxX.EoccJgsFoyJ/JcQNMShZHq', -- bcrypt hash for 'AdminBPS1200!'
    'internal'
)
ON CONFLICT (email) DO UPDATE 
SET username = 'admin',
    password_hash = '$2a$10$CtL0O5Pagh8WXov1KJIqVOmf7w0WxX.EoccJgsFoyJ/JcQNMShZHq';

-- Ensure admin role is assigned to the default admin user
INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u, roles r
WHERE u.username = 'admin' AND r.nama = 'admin'
ON CONFLICT DO NOTHING;
