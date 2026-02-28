-- name: UpdateCreditLimit :exec
UPDATE account_balances
SET credit_limit = $1,
    updated_at = now()
WHERE account_id = $2;