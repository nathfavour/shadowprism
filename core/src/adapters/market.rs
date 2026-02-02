use reqwest::Client;
use serde::Deserialize;
use std::time::{Duration, Instant};
use std::sync::Mutex;

#[derive(Deserialize, Debug)]
struct EncryptTradePrice {
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
                .timeout(Duration::from_secs(2))
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

        // 2. Fetch from Encrypt.trade (Mocked for Hackathon)
        // In production: https://api.encrypt.trade/v1/price/SOL-USD
        let price = match self.fetch_price_from_api().await {
            Ok(p) => p,
            Err(_) => {
                // Fallback to a static price if API is down
                println!("⚠️  Encrypt.trade API unreachable, using fallback price.");
                142.65 
            }
        };

        // 3. Update Cache
        let mut cache = self.cache.lock().unwrap();
        *cache = Some((price, Instant::now()));

        price
    }

    async fn fetch_price_from_api(&self) -> Result<f64, String> {
        // Mocking the API response for the hackathon demo
        // This demonstrates the integration point
        Ok(142.65)
    }

    pub fn format_usd(&self, lamports: u64, sol_price: f64) -> String {
        let sol = lamports as f64 / 1_000_000_000.0;
        let usd = sol * sol_price;
        format!("${:.2}", usd)
    }
}
