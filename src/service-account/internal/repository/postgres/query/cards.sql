-- name: CreateCard :one
INSERT INTO cards (id,
                   account_id,
                   pan_mask,
                   expiry_date,
                   is_virtual,
                   status)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: GetCardByID :one
SELECT
    id,
    account_id,
    pan_mask,
    expiry_date,
    is_virtual,
    status,
    created_at
FROM cards
WHERE id = $1 LIMIT 1;