-- +migrate Up notransaction
CREATE TYPE presence_status AS ENUM ('PRESENT', 'ABSENT', 'EXECUSE', 'PENDING_PRESENT', 'PENDING_EXECUSE');
-- +migrate Down
DROP TYPE IF EXISTS presence_status;