-- +goose Up
ALTER TABLE files
ADD COLUMN message_id UUID REFERENCES messages(id) ON DELETE SET NULL;    

-- +goose Down
ALTER TABLE files
DROP COLUMN message_id;
