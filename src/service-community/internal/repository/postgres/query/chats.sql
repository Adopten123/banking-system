-- name: CreateMessage :one
-- Создает новое сообщение в чате (обычное или перевод)
INSERT INTO messages (chat_id,
                      sender_id,
                      reply_to_message_id,
                      content,
                      media_attachments,
                      is_transfer,
                      transfer_amount,
                      transfer_currency,
                      idempotency_key,
                      transfer_transaction_id,
                      transfer_status)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING *;

-- name: GetChatMessages :many
-- Получает историю сообщений конкретного чата с пагинацией (от новых к старым)
SELECT id,
       chat_id,
       sender_id,
       reply_to_message_id,
       content,
       media_attachments,
       is_transfer,
       transfer_amount,
       transfer_currency,
       idempotency_key,
       transfer_transaction_id,
       transfer_status,
       is_edited,
       created_at,
       deleted_at
FROM messages
WHERE chat_id = $1
  AND deleted_at IS NULL
ORDER BY created_at DESC LIMIT $2
OFFSET $3;

-- name: CreateChat :one
INSERT INTO chats (type_id, title, avatar_url)
VALUES ($1, $2, $3)
    RETURNING *;

-- name: AddChatMember :exec
INSERT INTO chat_members (chat_id, user_id, role)
VALUES ($1, $2, $3);

-- name: GetChatMemberIDs :many
SELECT user_id FROM chat_members WHERE chat_id = $1;