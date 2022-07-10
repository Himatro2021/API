-- +migrate Up notransaction
CREATE TABLE IF NOT EXISTS user_invitations (
    id BIGINT PRIMARY KEY,
    email TEXT NOT NULL,
    invitation_code TEXT NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS user_invitations;