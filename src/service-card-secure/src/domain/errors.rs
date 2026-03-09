use std::fmt;

#[derive(Debug)]
pub enum VaultError {
    CardNotFound,
    EncryptionFailed(String),
    DecryptionFailed(String),
    HashFailed(String),
    InvalidPin,
    CardBlocked,
    InternalError(String),
}

impl fmt::Display for VaultError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            VaultError::CardNotFound => write!(f, "Card not found"),
            VaultError::EncryptionFailed(msg) => write!(f, "Encryption error: {}", msg),
            VaultError::DecryptionFailed(msg) => write!(f, "Decryption error: {}", msg),
            VaultError::HashFailed(msg) => write!(f, "Hash Error: {}", msg),
            VaultError::InvalidPin => write!(f, "Incorrect PIN-code"),
            VaultError::CardBlocked => write!(f, "Transaction declined: card blocked"),
            VaultError::InternalError(msg) => write!(f, "Internal error: {}", msg),
        }
    }
}

impl std::error::Error for VaultError {}