use crate::domain::errors::VaultError;
use crate::domain::models::{CardStatus, VaultCard};
use crate::domain::ports::CardRepository;
use sqlx::{PgPool, Row};
use uuid::Uuid;

pub struct PostgresCardRepository {
    pool: PgPool,
}

impl PostgresCardRepository {
    pub fn new(pool: PgPool) -> Self {
        Self { pool }
    }

    fn map_row_to_card(row: sqlx::postgres::PgRow) -> Result<VaultCard, VaultError> {
        let status_str: String = row.try_get("status")
            .map_err(|e| VaultError::InternalError(e.to_string()))?;

        let status = match status_str.as_str() {
            "ACTIVE" => CardStatus::Active,
            "BLOCKED" => CardStatus::Blocked,
            _ => return Err(VaultError::InternalError(format!("Unknown status: {}", status_str))),
        };

        Ok(VaultCard {
            token_id: row.try_get("token_id").unwrap_or_default(),
            pan_hash: row.try_get("pan_hash").unwrap_or_default(),
            encrypted_pan: row.try_get("encrypted_pan").unwrap_or_default(),
            encrypted_cvv: row.try_get("encrypted_cvv").unwrap_or_default(),
            expiry_month: row.try_get("expiry_month").unwrap_or_default(),
            expiry_year: row.try_get("expiry_year").unwrap_or_default(),
            status,
            pin_hash: row.try_get("pin_hash").unwrap_or_default(),
            created_at: row.try_get("created_at").unwrap_or_default(),
            updated_at: row.try_get("updated_at").unwrap_or_default(),
        })
    }
}

#[async_trait::async_trait]
impl CardRepository for PostgresCardRepository {
    async fn insert(&self, card: VaultCard) -> Result<(), VaultError> {
        let query = r#"
            INSERT INTO vault_cards (
                token_id, pan_hash, encrypted_pan, encrypted_cvv,
                expiry_month, expiry_year, status, pin_hash, created_at, updated_at
            ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        "#;

        sqlx::query(query)
            .bind(card.token_id)
            .bind(card.pan_hash)
            .bind(card.encrypted_pan)
            .bind(card.encrypted_cvv)
            .bind(card.expiry_month)
            .bind(card.expiry_year)
            .bind(card.status.as_str())
            .bind(card.pin_hash)
            .bind(card.created_at)
            .bind(card.updated_at)
            .execute(&self.pool)
            .await
            .map_err(|e| VaultError::InternalError(format!("Error inserting into database: {}", e)))?;

        Ok(())
    }

    async fn find_by_token(&self, token_id: &Uuid) -> Result<Option<VaultCard>, VaultError> {
        let query = "SELECT * FROM vault_cards WHERE token_id = $1";

        let result = sqlx::query(query)
            .bind(token_id)
            .fetch_optional(&self.pool)
            .await
            .map_err(|e| VaultError::InternalError(format!("Token search error: {}", e)))?;

        match result {
            Some(row) => Ok(Some(Self::map_row_to_card(row)?)),
            None => Ok(None),
        }
    }

    async fn find_by_pan_hash(&self, pan_hash: &str) -> Result<Option<VaultCard>, VaultError> {
        let query = "SELECT * FROM vault_cards WHERE pan_hash = $1";

        let result = sqlx::query(query)
            .bind(pan_hash)
            .fetch_optional(&self.pool)
            .await
            .map_err(|e| VaultError::InternalError(format!("Hash lookup error: {}", e)))?;

        match result {
            Some(row) => Ok(Some(Self::map_row_to_card(row)?)),
            None => Ok(None),
        }
    }

    async fn update_status(&self, token_id: &Uuid, status: CardStatus) -> Result<(), VaultError> {
        let query = "UPDATE vault_cards SET status = $1, updated_at = NOW() WHERE token_id = $2";

        let result = sqlx::query(query)
            .bind(status.as_str())
            .bind(token_id)
            .execute(&self.pool)
            .await
            .map_err(|e| VaultError::InternalError(format!("Status update error: {}", e)))?;

        if result.rows_affected() == 0 {
            return Err(VaultError::CardNotFound);
        }

        Ok(())
    }

    async fn delete(&self, token_id: &Uuid) -> Result<(), VaultError> {
        let query = "DELETE FROM vault_cards WHERE token_id = $1";

        let result = sqlx::query(query)
            .bind(token_id)
            .execute(&self.pool)
            .await
            .map_err(|e| VaultError::InternalError(format!("Uninstall error: {}", e)))?;

        if result.rows_affected() == 0 {
            return Err(VaultError::CardNotFound);
        }

        Ok(())
    }

    async fn update_pin(&self, token_id: &Uuid, pin_hash: &str) -> Result<(), VaultError> {
        let query = "UPDATE vault_cards SET pin_hash = $1, updated_at = NOW() WHERE token_id = $2";

        let result = sqlx::query(query)
            .bind(pin_hash)
            .bind(token_id)
            .execute(&self.pool)
            .await
            .map_err(|e| VaultError::InternalError(format!("PIN update error: {}", e)))?;

        if result.rows_affected() == 0 {
            return Err(VaultError::CardNotFound);
        }

        Ok(())
    }
}