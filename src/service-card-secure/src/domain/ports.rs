use crate::domain::errors::VaultError;
use crate::domain::models::{CardStatus, VaultCard};
use uuid::Uuid;

#[async_trait::async_trait]
pub trait CardRepository: Send + Sync {
    async fn insert(&self, card: VaultCard) -> Result<(), VaultError>;
    async fn find_by_token(&self, token_id: &Uuid) -> Result<Option<VaultCard>, VaultError>;
    async fn find_by_pan_hash(&self, pan_hash: &str) -> Result<Option<VaultCard>, VaultError>;
    async fn update_status(&self, token_id: &Uuid, status: CardStatus) -> Result<(), VaultError>;
    async fn delete(&self, token_id: &Uuid) -> Result<(), VaultError>;
    async fn update_pin(&self, token_id: &Uuid, pin_hash: &str) -> Result<(), VaultError>;
}

pub trait CryptoProvider: Send + Sync {
    fn encrypt(&self, plain_text: &str) -> Result<Vec<u8>, VaultError>;
    fn decrypt(&self, cipher_text: &[u8]) -> Result<String, VaultError>;
    fn hash_pan(&self, pan: &str) -> Result<String, VaultError>;
    fn hash_pin(&self, pin: &str) -> Result<String, VaultError>;
    fn verify_pin(&self, pin: &str, hash: &str) -> Result<bool, VaultError>;
}