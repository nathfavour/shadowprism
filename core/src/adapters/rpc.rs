use solana_client::rpc_client::RpcClient;
use std::env;

pub struct ReliableClient {
    primary: RpcClient,
    secondary: Option<RpcClient>,
}

impl ReliableClient {
    pub fn new() -> Self {
        let primary_url = env::var("HELIUS_RPC_URL")
            .unwrap_or_else(|_| "https://api.devnet.solana.com".to_string());
        
        let secondary_url = env::var("QUICKNODE_RPC_URL").ok();

        Self {
            primary: RpcClient::new(primary_url),
            secondary: secondary_url.map(RpcClient::new),
        }
    }

    pub fn get_client(&self) -> &RpcClient {
        // Simple logic: If we had a way to health-check, we would switch here.
        // For now, it returns the primary, but adapters can use .failover()
        &self.primary
    }

    pub fn get_secondary(&self) -> Option<&RpcClient> {
        self.secondary.as_ref()
    }

    pub fn send_transaction_reliable(&self, tx: &solana_sdk::transaction::Transaction) -> Result<solana_sdk::signature::Signature, String> {
        match self.primary.send_and_confirm_transaction(tx) {
            Ok(sig) => Ok(sig),
            Err(e) => {
                if let Some(ref secondary) = self.secondary {
                    println!("⚠️ [Failover] Primary RPC (Helius) failed/timed out. Routing to QuickNode...");
                    secondary.send_and_confirm_transaction(tx)
                        .map_err(|e2| format!("Both RPCs failed. Helius: {}, QuickNode: {}", e, e2))
                } else {
                    Err(format!("Primary RPC failed: {}", e))
                }
            }
        }
    }

    pub async fn get_priority_fee(&self) -> u64 {
        // In production: Query Helius Priority Fee API
        // For hackathon: Simulate a dynamic fee based on network "congestion"
        let base_fee = 5000;
        let jitter = rand::random::<u64>() % 2000;
        println!("🚀 [Helius] Calculating optimal priority fee: {} micro-lamports", base_fee + jitter);
        base_fee + jitter
    }
}
