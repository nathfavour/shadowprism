use std::sync::{Arc, Mutex};
use std::time::Duration;
use tokio::time::sleep;
use crate::db::TransactionStore;
use solana_client::rpc_client::RpcClient;
use solana_sdk::signature::Signature;
use std::str::FromStr;

pub struct Watchdog {
    store: Arc<Mutex<TransactionStore>>,
    rpc: RpcClient,
}

impl Watchdog {
    pub fn new(store: Arc<Mutex<TransactionStore>>) -> Self {
        let rpc_url = "https://api.devnet.solana.com".to_string();
        Self {
            store,
            rpc: RpcClient::new(rpc_url),
        }
    }

    pub async fn start(self) {
        println!("🐕 Transaction Watchdog started.");
        loop {
            self.check_pending_transactions().await;
            sleep(Duration::from_secs(30)).await;
        }
    }

    async fn check_pending_transactions(&self) {
        let pending_tasks = {
            let store = self.store.lock().unwrap();
            // We'll reuse list_transactions but filter locally for this example
            store.list_transactions().unwrap_or_default()
        };

        for task in pending_tasks {
            if task.status == "Pending" || task.status == "Broadcasted" {
                if let Some(hash) = task.tx_hash {
                    if let Ok(sig) = Signature::from_str(&hash) {
                        println!("🐕 Checking status for tx: {}", hash);
                        match self.rpc.get_signature_status(&sig) {
                            Ok(Some(Ok(_))) => {
                                println!("✅ Watchdog confirmed tx: {}", hash);
                                let store = self.store.lock().unwrap();
                                let _ = store.update_status(&task.id, "Confirmed", Some(&hash));
                            },
                            Ok(Some(Err(e))) => {
                                println!("❌ Watchdog found failed tx: {}", e);
                                let store = self.store.lock().unwrap();
                                let _ = store.update_status(&task.id, "Failed", Some(&hash));
                            },
                            _ => {} // Still pending or not found yet
                        }
                    }
                }
            }
        }
    }
}
