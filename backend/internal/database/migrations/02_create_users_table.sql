CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Insert default admin user (password should be hashed in production)
DELETE FROM users WHERE username='admin';
INSERT INTO users (username, password, role) VALUES ('admin', ADMIN_USER_PASSWORD_HASH, 'admin')
ON CONFLICT DO NOTHING;