use chrono::{DateTime, Utc};
use uuid::Uuid;

#[derive(Debug, Clone, PartialEq)]
pub enum CardStatus {
    Active,
    Blocked,
}

impl CardStatus {
    pub fn as_str(&self) -> &'static str {
        match self {
            CardStatus::Active => "ACTIVE",
            CardStatus::Blocked => "BLOCKED",
        }
    }
}

#[derive(Debug, Clone)]
pub struct VaultCard {
    pub token_id: Uuid,
    pub pan_hash: String,
    pub encrypted_pan: Vec<u8>,
    pub encrypted_cvv: Vec<u8>,
    pub expiry_month: i32,
    pub expiry_year: i32,
    pub status: CardStatus,
    pub pin_hash: Option<String>,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

#[derive(Debug, Clone)]
pub struct PlainCardDetails {
    pub pan: String,
    pub cvv: String,
    pub expiry_month: i32,
    pub expiry_year: i32,
}