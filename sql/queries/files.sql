-- name: CreateFile :exec
INSERT INTO files (
    id,
    chat_id,
    uploader_id,
    file_key,
    file_name,
    mime_type,
    size_bytes,
    message_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
);

-- name: DeleteFile :exec
DELETE FROM files
WHERE id = $1;

-- name: GetFileByID :one
SELECT 
    id, 
    chat_id, 
    uploader_id, 
    file_key, 
    file_name, 
    mime_type, 
    size_bytes, 
    updated_at
FROM files
WHERE id = $1;

-- name: LinkFileToMessage :exec
UPDATE files
SET message_id = $2
WHERE id = $1;