use reqwest::Client;
use std::time::{Duration, Instant};
use std::sync::Mutex;
use serde::Deserialize;

#[derive(Deserialize)]
struct JupiterPriceResponse {
    data: std::collections::HashMap<String, JupiterPriceData>,
}

#[derive(Deserialize)]
struct JupiterPriceData {
    price: f64,
}

pub struct MarketOracle {
    client: Client,
    cache: Mutex<Option<(f64, Instant)>>,
}

impl MarketOracle {
    pub fn new() -> Self {
        Self {
            client: Client::builder()
                .timeout(Duration::from_secs(5))
                .build()
                .unwrap(),
            cache: Mutex::new(None),
        }
    }

    pub async fn get_sol_price(&self) -> f64 {
        // 1. Check Cache (5 minute TTL)
        {
            let cache = self.cache.lock().unwrap();
            if let Some((price, timestamp)) = *cache {
                if timestamp.elapsed() < Duration::from_secs(300) {
                    return price;
                }
            }
        }

        // 2. Fetch from Jupiter Price API (Real Data)
        let price = match self.fetch_price_from_jup().await {
            Ok(p) => p,
            Err(e) => {
                println!("⚠️  Market API error: {}. Using fallback price.", e);
                142.65 
            }
        };

        // 3. Update Cache
        let mut cache = self.cache.lock().unwrap();
        *cache = Some((price, Instant::now()));

        price
    }

    async fn fetch_price_from_jup(&self) -> Result<f64, String> {
        let url = "https://api.jup.ag/price/v2?ids=So11111111111111111111111111111111111111112";
        let resp = self.client.get(url)
            .send()
            .await
            .map_err(|e| e.to_string())?;
        
        let data: JupiterPriceResponse = resp.json().await.map_err(|e| e.to_string())?;
        
        data.data.get("So11111111111111111111111111111111111111112")
            .map(|d| d.price)
            .ok_or_else(|| "SOL price not found in response".to_string())
    }

    pub fn format_usd(&self, lamports: u64, sol_price: f64) -> String {
        let sol = lamports as f64 / 1_000_000_000.0;
        let usd = sol * sol_price;
        format!("${:.2}", usd)
    }
}
