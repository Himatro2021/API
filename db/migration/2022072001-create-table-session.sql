-- +migrate Up notransaction
CREATE TABLE IF NOT EXISTS sessions (
    id BIGINT PRIMARY KEY,
    user_id BIGINT,
    access_token TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    access_token_expired_at TIMESTAMP NOT NULL,
    refresh_token_expired_at TIMESTAMP NOT NULL
);

ALTER TABLE sessions ADD FOREIGN KEY ("user_id") REFERENCES users("id");

-- +migrate Down
DROP TABLE IF EXISTS sessions;
