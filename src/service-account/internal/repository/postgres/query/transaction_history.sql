-- name: GetAccountTransactions :many
SELECT
    t.id AS transaction_id,
    t.category_id,
    t.status_id,
    t.description,
    t.created_at,
    p.amount::text AS amount_str,
    p.currency_code
FROM transactions t
JOIN postings p ON t.id = p.transaction_id
WHERE p.account_id = $1
    AND (sqlc.narg('start_date')::timestamp IS NULL OR t.created_at >= sqlc.narg('start_date'))
    AND (sqlc.narg('end_date')::timestamp IS NULL OR t.created_at <= sqlc.narg('end_date'))
ORDER BY t.created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');