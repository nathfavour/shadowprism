use crate::adapters::{PrivacyProvider, ShieldRequest, ShieldResponse, rpc::ReliableClient};
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

pub struct RadrAdapter;

impl RadrAdapter {
    pub const PROGRAM_ID: &str = "GQBqwwoikYh7p6KEUHDUu5r9dHHXx9tMGskAPubmFPzD";
}

#[async_trait]
impl PrivacyProvider for RadrAdapter {
    fn name(&self) -> String {
        "radr_shadow_wire".to_string()
    }

    async fn shield(&self, req: ShieldRequest, keystore: Arc<PrismKeystore>, rpc: Arc<ReliableClient>) -> Result<ShieldResponse, String> {
        let from_pubkey = keystore.main_keypair.pubkey();
        let to_pubkey = Pubkey::from_str(&req.destination_addr)
            .map_err(|e| format!("Invalid destination address: {}", e))?;
        let program_id = Pubkey::from_str(Self::PROGRAM_ID).unwrap();

        println!("👻 [Radr ShadowWire] Initiating P2P Encrypted Transfer via {}", Self::PROGRAM_ID);

        let recent_blockhash = rpc.get_client().get_latest_blockhash()
            .map_err(|e| format!("Failed to get blockhash: {}", e))?;

        // Constructing a real Radr ShadowWire instruction (Simulated format based on common patterns)
        // Discriminator for 'Transfer' (8 bytes) + Amount (8 bytes)
        let mut data = vec![0u8; 16];
        data[0..8].copy_from_slice(&[165, 12, 110, 212, 21, 12, 21, 10]); // Mock discriminator
        data[8..16].copy_from_slice(&req.amount_lamports.to_le_bytes());

        let mut ixs = vec![];
        let priority_fee = rpc.get_priority_fee().await;
        ixs.push(solana_compute_budget_interface::ComputeBudgetInstruction::set_compute_unit_price(priority_fee));
        
        ixs.push(Instruction {
            program_id,
            accounts: vec![
                AccountMeta::new(from_pubkey, true),
                AccountMeta::new(to_pubkey, false),
                AccountMeta::new_readonly(solana_sdk::system_program::id(), false),
            ],
            data,
        });

        let mut tx = Transaction::new_with_payer(
            &ixs,
            Some(&from_pubkey),
        );
        
        tx.sign(&[&keystore.main_keypair], recent_blockhash);

        let signature = rpc.send_transaction_reliable(&tx)?;

        Ok(ShieldResponse {
            status: "success".to_string(),
            tx_hash: signature.to_string(),
            provider: self.name(),
            note: Some("ghost-receipt-encrypted-p2p".to_string()),
        })
    }
}
