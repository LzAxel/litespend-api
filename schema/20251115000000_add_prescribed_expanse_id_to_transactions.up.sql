ALTER TABLE transactions ADD COLUMN prescribed_expanse_id BIGINT REFERENCES prescribed_expanses(id);

