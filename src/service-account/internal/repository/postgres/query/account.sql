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
