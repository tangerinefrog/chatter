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
UPDATE messages m
SET 
    read_at = now()
WHERE 
    m.chat_id = $2
    AND 
    m.user_id != $3
    AND 
    m.read_at IS NULL
    AND
    m.created_at <= (SELECT created_at FROM messages WHERE messages.id = $1);