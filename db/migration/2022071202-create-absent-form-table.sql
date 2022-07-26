-- +migrate Up notransaction
CREATE TABLE IF NOT EXISTS absent_forms (
    id BIGINT PRIMARY KEY,
    participant_group_id BIGINT NOT NULL,
    Title VARCHAR(255),
    start_at TIMESTAMP NOT NULL,
    finished_at TIMESTAMP NOT NULL,
    allow_update_by_attendee BOOLEAN DEFAULT FALSE,
    allow_create_confirmation BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP DEFAULT NULL,
    created_by BIGINT NOT NULL,
    updated_by BIGINT NOT NULL,
    deleted_by BIGINT DEFAULT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS absent_forms;
