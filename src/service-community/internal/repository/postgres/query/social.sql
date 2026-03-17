-- name: CreateReport :one
INSERT INTO reports (reporter_id, target_type, target_id, reason)
VALUES ($1, $2, $3, $4) RETURNING *;