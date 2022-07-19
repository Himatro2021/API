-- +migrate Up notransaction
ALTER TABLE absent_forms ADD FOREIGN KEY ("participant_group_id") REFERENCES group_members ("id");
ALTER TABLE absent_forms ADD FOREIGN KEY ("created_by") REFERENCES users ("id");
ALTER TABLE absent_forms ADD FOREIGN KEY ("updated_by") REFERENCES users ("id");
ALTER TABLE absent_forms ADD FOREIGN KEY ("deleted_by") REFERENCES users ("id");

-- +migrate Down