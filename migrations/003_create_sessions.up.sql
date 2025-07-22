CREATE TABLE sessions (
    token TEXT PRIMARY KEY,
    user_login TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);