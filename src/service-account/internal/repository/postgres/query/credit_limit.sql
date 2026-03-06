-- name: GetAccountInfoForLimitUpdate :one
SELECT
    a.id,
    a.currency_code,
    ab.credit_limit::text AS credit_limit_str
FROM accounts a
         JOIN account_balances ab ON a.id = ab.account_id
WHERE a.public_id = sqlc.arg('public_id');

-- name: UpdateAccountCreditLimit :exec
UPDATE account_balances
SET
    credit_limit = sqlc.arg('new_limit'),
    updated_at = now()
WHERE account_id = sqlc.arg('account_id');

-- name: GetCreditLimit :one
SELECT credit_limit::text
FROM account_balances
WHERE account_id = $1;