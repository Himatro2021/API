-- +migrate Up notransaction
ALTER TABLE user_invitations ADD COLUMN IF NOT EXISTS role user_roles DEFAULT 'MEMBER'::user_roles;

-- +migrate Down
ALTER TABLE user_invitations DROP COLUMN IF EXISTS role;