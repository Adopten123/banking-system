-- name: CreateRecurringPayment :one
INSERT INTO recurring_payments (
    id,
    source_type_id,
    source_id,
    destination_type_id,
    destination_id,
    amount,
    currency_code,
    category_id,
    cron_expression,
    next_execution_time,
    description
) VALUES (
             $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
         ) RETURNING *;

-- name: GetRecurringPaymentByID :one
SELECT *
FROM recurring_payments
WHERE id = $1;

-- name: UpdateRecurringPaymentStatus :exec
UPDATE recurring_payments
SET is_active = $2,
    updated_at = NOW()
WHERE id = $1;

-- name: GetDueRecurringPayments :many
SELECT *
FROM recurring_payments
WHERE is_active = true
  AND next_execution_time <= NOW()
    FOR UPDATE SKIP LOCKED
LIMIT $1;

-- name: UpdateRecurringPaymentNextRun :exec
UPDATE recurring_payments
SET next_execution_time = $2,
    updated_at = NOW()
WHERE id = $1;