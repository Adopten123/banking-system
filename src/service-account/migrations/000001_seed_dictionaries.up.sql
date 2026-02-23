-- 1. Testing user
INSERT INTO "users" ("id", "username")
VALUES ('123e4567-e89b-12d3-a456-426614174000', 'test_user') ON CONFLICT ("id") DO NOTHING;

-- 2. Add account types
INSERT INTO "account_types" ("id", "name")
VALUES (1, 'current'),
       (2, 'savings'),
       (3, 'credit') ON CONFLICT ("id") DO NOTHING;

-- 3. Add account statuses
INSERT INTO "account_statuses" ("id", "name")
VALUES (1, 'active'),
       (2, 'frozen'),
       (3, 'blocked'),
       (4, 'closed') ON CONFLICT ("id") DO NOTHING;

-- 4. Add currencies
INSERT INTO "currencies" ("code", "name", "exponent")
VALUES ('RUB', 'Russian Ruble', 2),
       ('USD', 'US Dollar', 2),
       ('EUR', 'Euro', 2) ON CONFLICT ("code") DO NOTHING;

-- 5. Add transaction categories
INSERT INTO "transaction_categories" ("id", "name")
VALUES (1, 'deposit'),
       (2, 'withdrawal'),
       (3, 'transfer') ON CONFLICT ("id") DO NOTHING;

-- 6. Add transaction statuses
INSERT INTO "transaction_statuses" ("id", "name")
VALUES (1, 'pending'),
       (2, 'posted'),
       (3, 'failed'),
       (4, 'rolled_back') ON CONFLICT ("id") DO NOTHING;