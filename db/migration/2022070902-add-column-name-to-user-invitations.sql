-- +migrate Up notransaction
ALTER TABLE user_invitations ADD COLUMN IF NOT EXISTS name text;

-- +migrate Down
ALTER TABLE user_invitations DROP COLUMN IF EXISTS name;