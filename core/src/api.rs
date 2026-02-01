use axum::{
    extract::State,
    http::StatusCode,
    Json,
};
use crate::adapters::{ShieldRequest, ShieldResponse};
use std::sync::Arc;

pub struct AppState {
    pub range: crate::adapters::range::RangeClient,
    pub providers: Vec<Box<dyn crate::adapters::PrivacyProvider>>,
}

pub async fn shield_handler(
    State(state): State<Arc<AppState>>,
    Json(payload): Json<ShieldRequest>,
) -> Result<Json<ShieldResponse>, (StatusCode, String)> {
    // 1. Compliance Check
    let risk_score = state.range.check_risk(&payload.destination_addr).await
        .map_err(|e| (StatusCode::INTERNAL_SERVER_ERROR, e))?;
    
    if risk_score > 80 {
        return Err((StatusCode::FORBIDDEN, "High risk destination address".to_string()));
    }

    // 2. Route Selection
    let provider = if payload.strategy.contains("p2p") {
        &state.providers[1] // Radr
    } else {
        &state.providers[0] // Privacy Cash
    };

    // 3. Execution
    let result = provider.shield(payload).await
        .map_err(|e| (StatusCode::INTERNAL_SERVER_ERROR, e))?;

    Ok(Json(result))
}