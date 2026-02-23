-- name: CreateTransaction :one
INSERT INTO transactions (id, category_id, status_id, description, external_details, idempotency_key)
VALUES ($1, $2, $3, $4, $5,
        $6) RETURNING id, category_id, status_id, description, external_details, idempotency_key, created_at, updated_at;

-- name: CreatePosting :one
INSERT INTO postings (transaction_id, account_id, amount, currency_code, exchange_rate)
VALUES ($1, $2, $3, $4, $5) RETURNING id, transaction_id, account_id, amount, currency_code, exchange_rate;

-- name: AddAccountBalance :one
UPDATE account_balances
SET balance    = balance + $1,
    updated_at = now()
WHERE account_id = $2 RETURNING account_id, balance, credit_limit, updated_at;