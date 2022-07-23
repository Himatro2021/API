-- +migrate Up notransaction
CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY,
    email text UNIQUE NOT NULL,
    password text,
    name text
);

-- +migrate Down
DROP TABLE IF EXISTS users;