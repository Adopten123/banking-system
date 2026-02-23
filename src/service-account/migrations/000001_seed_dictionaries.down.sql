DELETE FROM "transaction_statuses" WHERE "id" IN (1, 2, 3, 4);
DELETE FROM "transaction_categories" WHERE "id" IN (1, 2, 3);
DELETE FROM "currencies" WHERE "code" IN ('RUB', 'USD', 'EUR');
DELETE FROM "account_statuses" WHERE "id" IN (1, 2, 3, 4);
DELETE FROM "account_types" WHERE "id" IN (1, 2, 3);
DELETE FROM "users" WHERE "id" = '123e4567-e89b-12d3-a456-426614174000';
