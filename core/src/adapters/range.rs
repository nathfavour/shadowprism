use reqwest::Client;
use serde::Deserialize;
use std::collections::HashMap;
use std::sync::Mutex;
use std::time::{Duration, Instant};

#[derive(Deserialize)]
struct RangeResponse {
    score: u8,
}

pub struct RangeClient {
    client: Client,
    api_key: String,
    cache: Mutex<HashMap<String, (u8, Instant)>>,
}

impl RangeClient {
    pub fn new() -> Self {
        let api_key = std::env::var("RANGE_API_KEY").unwrap_or_else(|_| "dev_default_key".to_string());
        Self {
            client: Client::builder()
                .timeout(Duration::from_secs(3))
                .build()
                .unwrap(),
            api_key,
            cache: Mutex::new(HashMap::new()),
        }
    }

    pub async fn check_risk(&self, address: &str) -> Result<u8, String> {
        // 1. Check Cache (1 hour TTL)
        {
            let cache = self.cache.lock().unwrap();
            if let Some((score, timestamp)) = cache.get(address) {
                if timestamp.elapsed() < Duration::from_secs(3600) {
                    return Ok(*score);
                }
            }
        }

        // 2. Real API Call
        // In production: https://api.rangeprotocol.com/v1/score/{address}
        // For hackathon/dev: we simulate the call if the key is default
        if self.api_key == "dev_default_key" {
            let score = if address.starts_with("BAD") { 99 } else { 10 };
            return Ok(score);
        }

        let url = format!("https://api.rangeprotocol.com/v1/score/{}", address);
        let resp = self.client
            .get(url)
            .header("X-API-KEY", &self.api_key)
            .send()
            .await
            .map_err(|e| format!("Range API error: {}", e))?;

        if resp.status().is_success() {
            let data: RangeResponse = resp.json().await.map_err(|e| e.to_string())?;
            
            // 3. Update Cache
            let mut cache = self.cache.lock().unwrap();
            cache.insert(address.to_string(), (data.score, Instant::now()));
            
            Ok(data.score)
        } else {
            Err(format!("Range API returned error: {}", resp.status()))
        }
    }
}