use async_trait::async_trait;
use serde::{Deserialize, Serialize};
use std::sync::Arc;
use crate::keystore::PrismKeystore;

#[derive(Debug, Serialize, Deserialize)]
pub struct ShieldRequest {
    pub amount_lamports: u64,
    pub destination_addr: String,
    pub strategy: String,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct ShieldResponse {
    pub status: String,
    pub tx_hash: String,
    pub provider: String,
    pub note: Option<String>,
}

#[async_trait]
pub trait PrivacyProvider: Send + Sync {
    async fn shield(&self, req: ShieldRequest, keystore: Arc<PrismKeystore>) -> Result<ShieldResponse, String>;
    fn name(&self) -> String;
}

pub mod privacy_cash;
pub mod radr;
pub mod range;
pub mod rpc;
