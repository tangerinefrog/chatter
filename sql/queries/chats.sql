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
    (
        SELECT json_agg(row_to_json(t))
        FROM (
            SELECT
                u.id,
                u.username
            FROM chats_users cu
            INNER JOIN users u on cu.user_id = u.id
            where cu.chat_id = c.id
        ) t
    ) chat_users_json
FROM chats c
INNER JOIN chats_users cu ON cu.chat_id = c.id
WHERE 
    cu.user_id = $1;
