-- name: CreateMessage :one
INSERT INTO messages (
    chat_id,
    user_id,
    content
)
VALUES (
    $1,
    $2,
    $3
)
RETURNING id;

-- name: ListTopNMessages :many
SELECT
    m.id,
    m.user_id,
    m.content,
    m.created_at,
    m.read_at
FROM messages m
WHERE m.chat_id = $1
ORDER BY m.created_at DESC
LIMIT $2 OFFSET $3;

-- name: MarkMessagesAsRead :exec
UPDATE messages
SET 
    read_at = now()
WHERE 
    id <= $1
    AND
    chat_id = $2
    AND 
    user_id != $3
    AND 
    read_at IS NULL;