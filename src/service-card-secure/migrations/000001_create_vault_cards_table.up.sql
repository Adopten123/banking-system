CREATE TABLE vault_cards
(
    token_id      UUID PRIMARY KEY,
    pan_hash      VARCHAR(64) UNIQUE NOT NULL,
    encrypted_pan BYTEA              NOT NULL,
    encrypted_cvv BYTEA              NOT NULL,
    expiry_month  INT                NOT NULL,
    expiry_year   INT                NOT NULL,
    status        VARCHAR(20) DEFAULT 'ACTIVE',
    pin_hash      VARCHAR(255),
    created_at    TIMESTAMPTZ DEFAULT NOW(),
    updated_at    TIMESTAMPTZ DEFAULT NOW()
);