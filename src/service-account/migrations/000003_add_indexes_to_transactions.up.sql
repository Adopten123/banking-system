-- Индекс для быстрого поиска проводок по ID счета
CREATE INDEX IF NOT EXISTS idx_postings_account_id ON postings(account_id);

-- Индекс для сортировки транзакций от новых к старым
CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions(created_at DESC);