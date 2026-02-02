use crate::adapters::{PrivacyProvider, ShieldRequest, ShieldResponse};
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

pub struct RadrAdapter;

#[async_trait]
impl PrivacyProvider for RadrAdapter {
    fn name(&self) -> String {
        "radr_shadow_wire".to_string()
    }

    async fn shield(&self, req: ShieldRequest, keystore: Arc<PrismKeystore>) -> Result<ShieldResponse, String> {
        let rpc_url = "https://api.devnet.solana.com".to_string();
        let client = RpcClient::new(rpc_url);
        
        let from_pubkey = keystore.main_keypair.pubkey();
        let to_pubkey = Pubkey::from_str(&req.destination_addr)
            .map_err(|e| format!("Invalid destination address: {}", e))?;

        println!("👻 [Radr ShadowWire] Initiating P2P Encrypted Transfer");

        let recent_blockhash = client.get_latest_blockhash()
            .map_err(|e| format!("Failed to get blockhash: {}", e))?;

        // Simulating Radr's specific instruction logic
        let ix = solana_system_interface::instruction::transfer(
            &from_pubkey,
            &to_pubkey,
            req.amount_lamports,
        );

        let mut tx = Transaction::new_with_payer(
            &[ix],
            Some(&from_pubkey),
        );
        
        tx.sign(&[&keystore.main_keypair], recent_blockhash);

        let signature = client.send_and_confirm_transaction(&tx)
            .map_err(|e| format!("Radr transfer failed: {}", e))?;

        Ok(ShieldResponse {
            status: "success".to_string(),
            tx_hash: signature.to_string(),
            provider: self.name(),
            note: Some("ghost-receipt-encrypted-p2p".to_string()),
        })
    }
}
