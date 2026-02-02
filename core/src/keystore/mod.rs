use solana_sdk::signature::{Keypair, read_keypair_file};
use std::path::{Path, PathBuf};
use anyhow::{Result, Context, anyhow};
use aes_gcm::{
    aead::{Aead, KeyInit},
    Aes256Gcm, Nonce
};
use pbkdf2::pbkdf2_hmac;
use sha2::Sha256;
use rand::{RngCore, thread_rng};
use base64::{Engine as _, engine::general_purpose::STANDARD as BASE64};
use std::fs;

pub struct PrismKeystore {
    pub main_keypair: Keypair,
}

impl PrismKeystore {
    /// Loads an encrypted keypair from disk or generates a new one if it doesn't exist.
    pub fn load_or_create(path: &Path, password: &str) -> Result<Self> {
        if path.exists() {
            Self::load(path, password)
        } else {
            let keypair = Keypair::new();
            Self::save(path, &keypair, password)?;
            Ok(Self { main_keypair: keypair })
        }
    }

    pub fn load(path: &Path, password: &str) -> Result<Self> {
        let content = fs::read_to_string(path).context("Failed to read keystore file")?;
        let decoded = BASE64.decode(content.trim()).map_err(|e| anyhow!("Invalid base64: {}", e))?;
        
        if decoded.len() < 32 {
            return Err(anyhow!("Keystore file too short"));
        }

        let salt = &decoded[0..16];
        let nonce_bytes = &decoded[16..28];
        let ciphertext = &decoded[28..];

        let mut key = [0u8; 32];
        pbkdf2_hmac::<Sha256>(password.as_bytes(), salt, 100_000, &mut key);

        let cipher = Aes256Gcm::new_from_slice(&key).map_err(|e| anyhow!("Cipher init error: {}", e))?;
        let nonce = Nonce::from_slice(nonce_bytes);
        
        let plaintext = cipher.decrypt(nonce, ciphertext)
            .map_err(|_| anyhow!("Decryption failed - wrong passphrase?"))?;

        let keypair = Keypair::from_bytes(&plaintext)
            .map_err(|e| anyhow!("Invalid keypair data: {}", e))?;

        Ok(Self { main_keypair: keypair })
    }

    pub fn save(path: &Path, keypair: &Keypair, password: &str) -> Result<()> {
        let mut salt = [0u8; 16];
        let mut nonce_bytes = [0u8; 12];
        thread_rng().fill_bytes(&mut salt);
        thread_rng().fill_bytes(&mut nonce_bytes);

        let mut key = [0u8; 32];
        pbkdf2_hmac::<Sha256>(password.as_bytes(), &salt, 100_000, &mut key);

        let cipher = Aes256Gcm::new_from_slice(&key).map_err(|e| anyhow!("Cipher init error: {}", e))?;
        let nonce = Nonce::from_slice(&nonce_bytes);

        let ciphertext = cipher.encrypt(nonce, keypair.to_bytes().as_ref())
            .map_err(|e| anyhow!("Encryption failed: {}", e))?;

        let mut output = Vec::with_capacity(16 + 12 + ciphertext.len());
        output.extend_from_slice(&salt);
        output.extend_from_slice(&nonce_bytes);
        output.extend_from_slice(&ciphertext);

        let encoded = BASE64.encode(output);
        fs::write(path, encoded).context("Failed to write keystore file")?;

        Ok(())
    }

    pub fn pubkey(&self) -> solana_sdk::pubkey::Pubkey {
        use solana_sdk::signer::Signer;
        self.main_keypair.pubkey()
    }
}
