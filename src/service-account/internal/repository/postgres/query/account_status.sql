-- name: UpdateAccountStatus :exec
UPDATE accounts
SET status_id = $1,
    updated_at = now()
WHERE id = $2;

-- name: GetBalanceForUpdate :one
SELECT balance, credit_limit
FROM account_balances
WHERE account_id = $1 FOR UPDATE;