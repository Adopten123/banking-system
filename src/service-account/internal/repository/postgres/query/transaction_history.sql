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
ORDER BY t.created_at DESC;