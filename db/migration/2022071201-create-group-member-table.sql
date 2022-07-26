-- +migrate Up notransaction
CREATE TABLE IF NOT EXISTS group_members (
    id BIGINT PRIMARY KEY,
    name TEXT NOT NULL,
    members BIGINT[] 
);

-- +migrate Down
DROP TABLE IF EXISTS group_members;