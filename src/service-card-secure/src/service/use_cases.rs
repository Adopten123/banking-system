use std::sync::Arc;
use chrono::Utc;
use uuid::Uuid;

use crate::domain::errors::VaultError;
use crate::domain::models::{CardStatus, PlainCardDetails, VaultCard};
use crate::domain::ports::{CardRepository, CryptoProvider};
use crate::infrastructure::generator::CardGenerator;

pub struct IssueCardResult {
    pub token_id: Uuid,
    pub pan_mask: String,
    pub expiry_month: i32,
    pub expiry_year: i32,
}

pub struct CardVaultUseCase {
    repo: Arc<dyn CardRepository>,
    crypto: Arc<dyn CryptoProvider>,
}

impl CardVaultUseCase {
    pub fn new(repo: Arc<dyn CardRepository>, crypto: Arc<dyn CryptoProvider>) -> Self {
        Self { repo, crypto }
    }

    fn get_prefix(payment_system: &str) -> &'static str {
        match payment_system.to_lowercase().as_str() {
            "visa" => "4276",
            "mastercard" => "5469",
            "mir" => "2200",
            _ => "4000",
        }
    }

    pub async fn issue_card(
        &self,
        payment_system: &str,
        _is_virtual: bool
    ) -> Result<IssueCardResult, VaultError> {

        let prefix = Self::get_prefix(payment_system);
        let pan = CardGenerator::generate_pan(prefix);
        let cvv = CardGenerator::generate_cvv();
        let (expiry_month, expiry_year) = CardGenerator::generate_expiry();

        let pan_hash = self.crypto.hash_pan(&pan)?;
        let encrypted_pan = self.crypto.encrypt(&pan)?;
        let encrypted_cvv = self.crypto.encrypt(&cvv)?;

        let token_id = Uuid::new_v4();
        let now = Utc::now();

        let card = VaultCard {
            token_id,
            pan_hash,
            encrypted_pan,
            encrypted_cvv,
            expiry_month,
            expiry_year,
            status: CardStatus::Active,
            pin_hash: None,
            created_at: now,
            updated_at: now,
        };

        self.repo.insert(card).await?;

        Ok(IssueCardResult {
            token_id,
            pan_mask: CardGenerator::mask_pan(&pan),
            expiry_month,
            expiry_year,
        })
    }

    pub async fn get_card_details(&self, token_id: &Uuid) -> Result<PlainCardDetails, VaultError> {
        let card = self.repo.find_by_token(token_id).await?
            .ok_or(VaultError::CardNotFound)?;

        if card.status == CardStatus::Blocked {
            return Err(VaultError::CardBlocked);
        }

        let pan = self.crypto.decrypt(&card.encrypted_pan)?;
        let cvv = self.crypto.decrypt(&card.encrypted_cvv)?;

        Ok(PlainCardDetails {
            pan,
            cvv,
            expiry_month: card.expiry_month,
            expiry_year: card.expiry_year,
        })
    }

    pub async fn verify_card(
        &self,
        pan: &str,
        cvv: &str,
        exp_month: i32,
        exp_year: i32
    ) -> Result<(bool, Option<Uuid>), VaultError> {

        let pan_hash = self.crypto.hash_pan(pan)?;

        let card = match self.repo.find_by_pan_hash(&pan_hash).await? {
            Some(c) => c,
            None => return Ok((false, None)),
        };

        if card.status == CardStatus::Blocked {
            return Ok((false, Some(card.token_id)));
        }

        let decrypted_cvv = self.crypto.decrypt(&card.encrypted_cvv)?;

        let is_valid = decrypted_cvv == cvv
            && card.expiry_month == exp_month
            && card.expiry_year == exp_year;

        Ok((is_valid, Some(card.token_id)))
    }

    pub async fn update_status(&self, token_id: &Uuid, status_str: &str) -> Result<(), VaultError> {
        let status = match status_str.to_uppercase().as_str() {
            "ACTIVE" => CardStatus::Active,
            "BLOCKED" => CardStatus::Blocked,
            _ => return Err(VaultError::InternalError("Invalid status".into())),
        };

        self.repo.update_status(token_id, status).await
    }

    pub async fn delete_card(&self, token_id: &Uuid) -> Result<(), VaultError> {
        self.repo.delete(token_id).await
    }

    pub async fn set_pin(&self, token_id: &Uuid, pin: &str) -> Result<(), VaultError> {
        let _ = self.repo.find_by_token(token_id).await?
            .ok_or(VaultError::CardNotFound)?;

        let pin_hash = self.crypto.hash_pin(pin)?;
        self.repo.update_pin(token_id, &pin_hash).await
    }

    pub async fn verify_pin(&self, token_id: &Uuid, pin: &str) -> Result<bool, VaultError> {
        let card = self.repo.find_by_token(token_id).await?
            .ok_or(VaultError::CardNotFound)?;

        if card.status == CardStatus::Blocked {
            return Err(VaultError::CardBlocked);
        }

        match card.pin_hash {
            Some(hash) => {
                let is_valid = self.crypto.verify_pin(pin, &hash)?;
                Ok(is_valid)
            },
            None => Ok(false),
        }
    }
}