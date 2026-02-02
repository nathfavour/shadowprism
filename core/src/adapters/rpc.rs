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
}
