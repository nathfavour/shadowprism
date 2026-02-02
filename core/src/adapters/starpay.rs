use crate::adapters::{PaymentProvider, PayRequest, PayResponse, rpc::ReliableClient};
use crate::keystore::PrismKeystore;
use async_trait::async_trait;
use std::sync::Arc;
use solana_sdk::{
    transaction::Transaction,
    signer::Signer,
    pubkey::Pubkey,
};
use std::str::FromStr;
use uuid::Uuid;

pub struct StarpayAdapter;

#[async_trait]
impl PaymentProvider for StarpayAdapter {
    async fn pay(&self, req: PayRequest, keystore: Arc<PrismKeystore>, rpc: Arc<ReliableClient>) -> Result<PayResponse, String> {
        let from_pubkey = keystore.main_keypair.pubkey();
        let merchant_pubkey = Pubkey::from_str(&req.merchant_id)
            .map_err(|e| format!("Invalid merchant ID: {}", e))?;

        println!("💳 [Starpay] Processing payment of {} lamports to merchant {}", 
            req.amount_lamports, req.merchant_id);

        let recent_blockhash = rpc.get_client().get_latest_blockhash()
            .map_err(|e| format!("Failed to get blockhash: {}", e))?;

        // Simulating Starpay payment instruction
        let ix = solana_system_interface::instruction::transfer(
            &from_pubkey,
            &merchant_pubkey,
            req.amount_lamports,
        );

        let mut tx = Transaction::new_with_payer(
            &[ix],
            Some(&from_pubkey),
        );
        
        tx.sign(&[&keystore.main_keypair], recent_blockhash);

        let signature = rpc.send_transaction_reliable(&tx)?;

        let receipt_id = format!("STAR-{}", Uuid::new_v4());

        Ok(PayResponse {
            status: "success".to_string(),
            tx_hash: signature.to_string(),
            receipt_id,
        })
    }
}
