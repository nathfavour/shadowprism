use async_trait::async_trait;
use serde::{Deserialize, Serialize};
use std::sync::Arc;
use crate::keystore::PrismKeystore;

#[derive(Debug, Serialize, Deserialize)]
pub struct ShieldRequest {
    pub amount_lamports: u64,
    pub destination_addr: String,
    pub strategy: String,
    pub force: Option<bool>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct ShieldResponse {
    pub status: String,
    pub tx_hash: String,
    pub provider: String,
    pub note: Option<String>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct SwapRequest {
    pub amount_lamports: u64,
    pub from_token: String,
    pub to_token: String,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct SwapResponse {
    pub status: String,
    pub tx_hash: String,
    pub from_amount: u64,
    pub to_amount: u64,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct PayRequest {
    pub merchant_id: String,
    pub amount_lamports: u64,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct PayResponse {
    pub status: String,
    pub tx_hash: String,
    pub receipt_id: String,
}

#[async_trait]
pub trait PrivacyProvider: Send + Sync {
    async fn shield(&self, req: ShieldRequest, keystore: Arc<PrismKeystore>) -> Result<ShieldResponse, String>;
    fn name(&self) -> String;
}

#[async_trait]
pub trait SwapProvider: Send + Sync {
    async fn swap(&self, req: SwapRequest, keystore: Arc<PrismKeystore>) -> Result<SwapResponse, String>;
}

#[async_trait]
pub trait PaymentProvider: Send + Sync {
    async fn pay(&self, req: PayRequest, keystore: Arc<PrismKeystore>) -> Result<PayResponse, String>;
}

pub mod privacy_cash;
pub mod radr;
pub mod range;
pub mod rpc;
pub mod market;
pub mod silent_swap;
pub mod starpay;
