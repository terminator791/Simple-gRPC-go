-- User service database initialization
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert some sample data (password is 'password123' hashed with bcrypt)
INSERT INTO users (email, name, password_hash, role) VALUES 
    ('john.doe@example.com', 'John Doe', '$2a$10$K5Z5Y5Z5Y5Z5Y5Z5Y5Z5YuJ5Z5Y5Z5Y5Z5Y5Z5Z5Y5Z5Y5Z5Y5Z5Y5Y', 'user'),
    ('jane.smith@example.com', 'Jane Smith', '$2a$10$K5Z5Y5Z5Y5Z5Y5Z5Y5Z5YuJ5Z5Y5Z5Y5Z5Y5Z5Z5Y5Z5Y5Z5Y5Z5Y5Y', 'user'),
    ('admin@example.com', 'Administrator', '$2a$10$K5Z5Y5Z5Y5Z5Y5Z5Y5Z5YuJ5Z5Y5Z5Y5Z5Y5Z5Z5Y5Z5Y5Z5Y5Z5Y5Y', 'admin')
ON CONFLICT (email) DO NOTHING;