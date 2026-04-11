-- +goose Up
ALTER TABLE messages
ADD COLUMN read_at TIMESTAMPTZ NULL;

-- +goose Down
ALTER TABLE messages
DROP COLUMN read_at;
