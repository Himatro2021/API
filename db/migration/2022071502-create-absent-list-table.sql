-- +migrate Up notransaction
CREATE TABLE IF NOT EXISTS absent_lists (
    id BIGINT PRIMARY KEY,
    absent_form_id BIGINT NOT NULL,
    created_by BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    "status" presence_status NOT NULL,
    execuse_reason TEXT DEFAULT NULL
);

ALTER TABLE absent_lists ADD FOREIGN KEY ("absent_form_id") REFERENCES absent_forms("id");
ALTER TABLE absent_lists ADD FOREIGN KEY ("created_by") REFERENCES users("id");

-- +migrate Down
DROP TABLE IF EXISTS absent_lists;