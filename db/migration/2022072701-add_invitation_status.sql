-- +migrate Up notransaction
CREATE TYPE invitation_status AS ENUM ('SENT', 'FAILED', 'PENDING');

ALTER TABLE user_invitations ADD COLUMN IF NOT EXISTS status invitation_status NOT NULL;

-- +migrate Down
ALTER TABLE user_invitations DROP COLUMN IF EXISTS status;