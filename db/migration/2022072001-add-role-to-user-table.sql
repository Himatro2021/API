-- +migrate Up notransaction
CREATE TYPE user_roles AS ENUM ('MEMBER', 'ADMIN');

ALTER TABLE users ADD COLUMN IF NOT EXISTS role user_roles NOT NULL DEFAULT 'MEMBER'::user_roles;

-- +migrate Down
DROP TYPE user_roles;

ALTER TABLE users DROP COLUMN IF EXISTS role;