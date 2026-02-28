-- name: GetAccountForUpdate :one
SELECT a.id,
       a.status_id,
       a.currency_code,
       ab.balance,
       ab.credit_limit
FROM accounts a
JOIN account_balances ab ON a.id = ab.account_id
WHERE a.id = $1
FOR UPDATE;