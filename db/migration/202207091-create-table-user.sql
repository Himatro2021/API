-- +migrate Up notransaction
CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY,
    email text,
    password text,
    name text
);

-- +migrate Down
DROP TABLE IF EXISTS users;