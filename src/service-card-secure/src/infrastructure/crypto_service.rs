use crate::domain::errors::VaultError;
use crate::domain::ports::CryptoProvider;

use aes_gcm::{
    aead::{Aead, AeadCore, KeyInit, OsRng},
    Aes256Gcm, Key, Nonce,
};
use argon2::{
    password_hash::{
        rand_core::OsRng as ArgonOsRng, PasswordHash, PasswordHasher, PasswordVerifier, SaltString,
    },
    Argon2,
};
use sha2::{Digest, Sha256};

pub struct CryptoService {
    cipher: Aes256Gcm,
}

impl CryptoService {
    pub fn new(master_key: &str) -> Result<Self, VaultError> {
        let key_bytes = master_key.as_bytes();
        if key_bytes.len() != 32 {
            return Err(VaultError::InternalError(
                "The master key must be exactly 32 bytes long.".into(),
            ));
        }

        let key = Key::<Aes256Gcm>::from_slice(key_bytes);
        let cipher = Aes256Gcm::new(key);

        Ok(Self { cipher })
    }
}

impl CryptoProvider for CryptoService {
    fn encrypt(&self, plain_text: &str) -> Result<Vec<u8>, VaultError> {
        let nonce = Aes256Gcm::generate_nonce(&mut OsRng);

        let cipher_text = self
            .cipher
            .encrypt(&nonce, plain_text.as_bytes())
            .map_err(|e| VaultError::EncryptionFailed(e.to_string()))?;

        let mut encrypted_data = nonce.to_vec();
        encrypted_data.extend_from_slice(&cipher_text);

        Ok(encrypted_data)
    }

    fn decrypt(&self, encrypted_data: &[u8]) -> Result<String, VaultError> {
        if encrypted_data.len() < 28 {
            return Err(VaultError::DecryptionFailed(
                "Incorrect encrypted data length".into(),
            ));
        }

        let (nonce_bytes, cipher_text) = encrypted_data.split_at(12);
        let nonce = Nonce::from_slice(nonce_bytes);

        let plain_text_bytes = self
            .cipher
            .decrypt(nonce, cipher_text)
            .map_err(|e| VaultError::DecryptionFailed(e.to_string()))?;

        String::from_utf8(plain_text_bytes)
            .map_err(|_| VaultError::DecryptionFailed("The data is not valid UTF-8".into()))
    }

    fn hash_pan(&self, pan: &str) -> Result<String, VaultError> {
        let mut hasher = Sha256::new();
        hasher.update(pan.as_bytes());
        let result = hasher.finalize();
        Ok(hex::encode(result))
    }

    fn hash_pin(&self, pin: &str) -> Result<String, VaultError> {
        let salt = SaltString::generate(&mut ArgonOsRng);

        let argon2 = Argon2::default();
        let password_hash = argon2
            .hash_password(pin.as_bytes(), &salt)
            .map_err(|e| VaultError::HashFailed(e.to_string()))?;

        Ok(password_hash.to_string())
    }

    fn verify_pin(&self, pin: &str, hash: &str) -> Result<bool, VaultError> {
        let parsed_hash = PasswordHash::new(hash)
            .map_err(|_| VaultError::HashFailed("Incorrect hash format in the database".into()))?;

        let is_valid = Argon2::default()
            .verify_password(pin.as_bytes(), &parsed_hash)
            .is_ok();

        Ok(is_valid)
    }
}