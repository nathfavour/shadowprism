use solana_sdk::signature::{Keypair, read_keypair_file};
use std::path::PathBuf;
use anyhow::{Result, Context};

pub struct PrismKeystore {
    pub main_keypair: Keypair,
}

impl PrismKeystore {
    pub fn new() -> Result<Self> {
        // Try loading from standard Solana CLI location first
        let mut keypair_path = home::home_dir()
            .context("Could not find home directory")?;
        keypair_path.push(".config/solana/id.json");

        let keypair = if keypair_path.exists() {
            read_keypair_file(&keypair_path)
                .map_err(|e| anyhow::anyhow!("Failed to read solana keypair: {}", e))?
        } else {
            // Fallback: Generate a new one if not found (for dev)
            println!("⚠️ No Solana CLI keypair found. Generating a temporary one.");
            Keypair::new()
        };

        Ok(Self { main_keypair: keypair })
    }

    pub fn pubkey(&self) -> solana_sdk::pubkey::Pubkey {
        use solana_sdk::signer::Signer;
        self.main_keypair.pubkey()
    }
}