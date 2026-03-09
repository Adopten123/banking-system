-- name: CreateCard :one
INSERT INTO cards (id,
                   account_id,
                   pan_mask,
                   expiry_date,
                   is_virtual,
                   status)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;