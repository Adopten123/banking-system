-- name: CreateAccount :one
INSERT INTO accounts (
    public_id, user_id, type_id, status_id, currency_code, name
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetAccountByID :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: GetAccountByPublicID :one
SELECT id, public_id, user_id, type_id, status_id, currency_code, name, version, created_at, updated_at
FROM accounts
WHERE public_id = $1 LIMIT 1;

-- name: GetAccountBalanceByPublicID :one
SELECT ab.balance::text
FROM account_balances ab
JOIN accounts a ON a.id = ab.account_id
WHERE a.public_id = $1 LIMIT 1;

-- name: CreateAccountBalance :exec
INSERT INTO account_balances (account_id, balance, credit_limit, updated_at)
VALUES ($1, 0, 0, now());
