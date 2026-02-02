use crate::adapters::{PaymentProvider, PayRequest, PayResponse};
use crate::keystore::PrismKeystore;
use async_trait::async_trait;
use std::sync::Arc;
use solana_client::rpc_client::RpcClient;
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
    async fn pay(&self, req: PayRequest, keystore: Arc<PrismKeystore>) -> Result<PayResponse, String> {
        let rpc_url = "https://api.devnet.solana.com".to_string();
        let client = RpcClient::new(rpc_url);
        
        let from_pubkey = keystore.main_keypair.pubkey();
        let merchant_pubkey = Pubkey::from_str(&req.merchant_id)
            .map_err(|e| format!("Invalid merchant ID: {}", e))?;

        println!("💳 [Starpay] Processing payment of {} lamports to merchant {}", 
            req.amount_lamports, req.merchant_id);

        let recent_blockhash = client.get_latest_blockhash()
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

        let signature = client.send_and_confirm_transaction(&tx)
            .map_err(|e| format!("Payment failed: {}", e))?;

        let receipt_id = format!("STAR-{}", Uuid::new_v4());

        Ok(PayResponse {
            status: "success".to_string(),
            tx_hash: signature.to_string(),
            receipt_id,
        })
    }
}
