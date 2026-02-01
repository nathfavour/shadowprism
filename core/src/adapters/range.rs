pub struct RangeClient {
    client: reqwest::Client,
}

impl RangeClient {
    pub fn new() -> Self {
        Self {
            client: reqwest::Client::new(),
        }
    }

    pub async fn check_risk(&self, address: &str) -> Result<u8, String> {
        // Mocking Range Protocol API call
        // In production: GET https://api.rangeprotocol.com/v1/score/{address}
        
        if address.starts_with("BAD") {
            return Ok(99); // High risk
        }
        
        Ok(0) // Low risk
    }
}
