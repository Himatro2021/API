-- +migrate Up notransaction
ALTER TABLE user_invitations ADD COLUMN IF NOT EXISTS mail_service_id BIGINT DEFAULT NULL;

-- +migrate Down
ALTER TABLE user_invitations DROP COLUMN IF EXISTS mail_service_id;