use solana_client::rpc_client::RpcClient;
use std::env;
use reqwest::Client;
use serde_json::json;

pub struct ReliableClient {
    primary: RpcClient,
    secondary: Option<RpcClient>,
    http_client: Client,
    primary_url: String,
}

impl ReliableClient {
    pub fn new() -> Self {
        let primary_url = env::var("HELIUS_RPC_URL")
            .unwrap_or_else(|_| "https://api.devnet.solana.com".to_string());
        
        let secondary_url = env::var("QUICKNODE_RPC_URL").ok();

        Self {
            primary: RpcClient::new(primary_url.clone()),
            secondary: secondary_url.map(RpcClient::new),
            http_client: Client::new(),
            primary_url,
        }
    }

    pub fn get_client(&self) -> &RpcClient {
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
                    println!("⚠️ [Failover] Primary RPC (Helius) failed. Routing to QuickNode...");
                    secondary.send_and_confirm_transaction(tx)
                        .map_err(|e2| format!("Both RPCs failed. Helius: {}, QuickNode: {}", e, e2))
                } else {
                    Err(format!("Primary RPC failed: {}", e))
                }
            }
        }
    }

    pub async fn get_priority_fee(&self) -> u64 {
        // Real Helius Priority Fee API call
        if self.primary_url.contains("helius-rpc.com") || self.primary_url.contains("helius.xyz") {
            let body = json!({
                "jsonrpc": "2.0",
                "id": "priority-fee-estimate",
                "method": "getPriorityFeeEstimate",
                "params": [{
                    "accountKeys": ["JUP6LkbZbjS1jKKpphsRLSKE6t124vR9f8jP26CAtv6"], // Sample high-activity account
                    "options": {
                        "recommended": true
                    }
                }]
            });

            if let Ok(resp) = self.http_client.post(&self.primary_url)
                .json(&body)
                .send()
                .await {
                if let Ok(json) = resp.json::<serde_json::Value>().await {
                    if let Some(estimate) = json["result"]["priorityFeeEstimate"].as_f64() {
                        println!("🚀 [Helius] Real-time priority fee estimate: {} micro-lamports", estimate);
                        return estimate as u64;
                    }
                }
            }
        }

        // Fallback to standard Solana RPC if Helius fails or isn't used
        match self.primary.get_recent_prioritization_fees(&[]) {
            Ok(fees) => {
                let avg = if !fees.is_empty() {
                    fees.iter().map(|f| f.prioritization_fee).sum::<u64>() / fees.len() as u64
                } else {
                    5000
                };
                println!("🚀 [RPC] Using average prioritization fee: {} micro-lamports", avg);
                avg
            },
            Err(_) => 5000,
        }
    }
}
