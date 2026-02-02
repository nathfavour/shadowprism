use crate::adapters::{PaymentProvider, PayRequest, PayResponse, rpc::ReliableClient};
use crate::keystore::PrismKeystore;
use async_trait::async_trait;
use std::sync::Arc;
use solana_sdk::{
    transaction::Transaction,
    signer::Signer,
    pubkey::Pubkey,
    instruction::{Instruction, AccountMeta},
};
use std::str::FromStr;
use uuid::Uuid;

pub struct StarpayAdapter;

impl StarpayAdapter {
    pub const PROGRAM_ID: &str = "SPayX11111111111111111111111111111111111111";
}

#[async_trait]
impl PaymentProvider for StarpayAdapter {
    async fn pay(&self, req: PayRequest, keystore: Arc<PrismKeystore>, rpc: Arc<ReliableClient>) -> Result<PayResponse, String> {
        let from_pubkey = keystore.main_keypair.pubkey();
        let merchant_pubkey = Pubkey::from_str(&req.merchant_id)
            .map_err(|e| format!("Invalid merchant ID: {}", e))?;
        let program_id = Pubkey::from_str(Self::PROGRAM_ID).unwrap();

        println!("💳 [Starpay] Processing payment of {} lamports to merchant {} via {}", 
            req.amount_lamports, req.merchant_id, Self::PROGRAM_ID);

        let recent_blockhash = rpc.get_client().get_latest_blockhash()
            .map_err(|e| format!("Failed to get blockhash: {}", e))?;

        // Constructing a realistic Starpay payment instruction
        // Discriminator for 'Settle' (8 bytes) + Amount (8 bytes)
        let mut data = vec![0u8; 16];
        data[0..8].copy_from_slice(&[105, 12, 110, 212, 21, 12, 21, 10]); // Settle discriminator
        data[8..16].copy_from_slice(&req.amount_lamports.to_le_bytes());

        let ix = Instruction {
            program_id,
            accounts: vec![
                AccountMeta::new(from_pubkey, true),
                AccountMeta::new(merchant_pubkey, false),
                AccountMeta::new_readonly(Pubkey::from_str("11111111111111111111111111111111").unwrap(), false),
            ],
            data,
        };

        let mut tx = Transaction::new_with_payer(
            &[ix],
            Some(&from_pubkey),
        );
        
        tx.sign(&[&keystore.main_keypair], recent_blockhash);

        let signature = rpc.send_transaction_reliable(&tx)?;

        let receipt_id = format!("STAR-{}", Uuid::new_v4().to_string().to_uppercase());

        Ok(PayResponse {
            status: "success".to_string(),
            tx_hash: signature.to_string(),
            receipt_id,
        })
    }
}
