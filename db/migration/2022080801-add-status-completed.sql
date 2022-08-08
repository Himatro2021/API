-- +migrate Up notransaction
ALTER TYPE invitation_status ADD VALUE 'COMPLETED';

-- +migrate Down
