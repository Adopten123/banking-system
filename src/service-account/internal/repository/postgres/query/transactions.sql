-- name: CreateTransaction :one
INSERT INTO transactions (
    id,
    source_type_id,
    source_id,
    category_id,
    status_id,
    description,
    external_details,
    idempotency_key
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    RETURNING id, source_type_id, source_id, category_id, status_id, description, external_details, idempotency_key, created_at, updated_at;

-- name: CreatePosting :one
INSERT INTO postings (transaction_id, account_id, amount, currency_code, exchange_rate)
VALUES ($1, $2, $3, $4, $5) RETURNING id, transaction_id, account_id, amount, currency_code, exchange_rate;

-- name: AddAccountBalance :one
UPDATE account_balances
SET balance    = balance + $1,
    updated_at = now()
WHERE account_id = $2 RETURNING account_id, balance, credit_limit, updated_at;

-- name: GetAccountForWithdrawUpdate :one
SELECT a.id, a.status_id, a.currency_code, ab.balance::text, ab.credit_limit::text
FROM accounts a
JOIN account_balances ab ON a.id = ab.account_id
WHERE a.public_id = $1
FOR NO KEY UPDATE;

-- name: SubtractAccountBalance :exec
UPDATE account_balances
SET balance = balance - $1
WHERE account_id = $2;