-- name: CreateChat :one
INSERT INTO chats (
    type,
    name,
    created_by
)
VALUES (
    $1,
    $2,
    $3
)
RETURNING id;

-- name: CreateChatUser :exec
INSERT INTO chats_users (
    chat_id,
    user_id
) VALUES (
    $1, $2
);

-- name: ListUserChats :many
SELECT
    c.id,
    c.type,
    c.name,
    c.created_at,
    (
        SELECT json_agg(row_to_json(t))
        FROM (
            SELECT
                u.id,
                u.username
            FROM chats_users cu
            INNER JOIN users u on cu.user_id = u.id
            WHERE cu.chat_id = c.id
        ) t
    ) chat_users_json,
    (
        SELECT row_to_json(t)
        FROM (
            SELECT 
                m.content,
                m.created_at
            FROM messages m
            WHERE m.chat_id = c.id
            ORDER BY m.created_at DESC
            LIMIT 1
        ) t
    ) last_message_json,
    (
        SELECT COUNT(*)
        FROM messages m
        WHERE 
        m.chat_id = c.id
        AND
        m.read_at IS NULL
        AND
        m.user_id != $1
    ) unread_messages_count
FROM chats c
INNER JOIN chats_users cu ON cu.chat_id = c.id
WHERE 
    cu.user_id = $1;

-- name: GetDirectChatBetweenUsers :one
SELECT 
    c.id
FROM chats c
JOIN chats_users cu1 ON cu1.chat_id = c.id
JOIN chats_users cu2 ON cu2.chat_id = c.id
WHERE 
    c.type = 'direct'
    AND cu1.user_id = $1
    AND cu2.user_id = $2
LIMIT 1;

-- name: ListChatUsers :many
SELECT
    cu.user_id,
    u.username
FROM chats_users cu
INNER JOIN users u ON u.id = cu.user_id
WHERE cu.chat_id = $1;