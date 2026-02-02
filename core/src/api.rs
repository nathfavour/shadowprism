use axum::{
    extract::{State, Path},
    http::StatusCode,
    Json,
};
use crate::adapters::{
    ShieldRequest, ShieldResponse, 
    SwapRequest, SwapResponse, 
    PayRequest, PayResponse,
    market::MarketOracle,
    SwapProvider, PaymentProvider
};
use crate::db::{TransactionStore, TransactionRecord};
use std::sync::{Arc, Mutex};
use serde_json::json;

pub struct AppState {
    pub rpc: Arc<crate::adapters::rpc::ReliableClient>,
    pub range: crate::adapters::range::RangeClient,
    pub market: MarketOracle,
    pub providers: Vec<Box<dyn crate::adapters::PrivacyProvider>>,
    pub swap_provider: Box<dyn SwapProvider>,
    pub pay_provider: Box<dyn PaymentProvider>,
    pub db: Arc<Mutex<TransactionStore>>,
    pub keystore: Arc<crate::keystore::PrismKeystore>,
}

pub async fn shield_handler(
    State(state): State<Arc<AppState>>,
    Json(payload): Json<ShieldRequest>,
) -> Result<Json<ShieldResponse>, (StatusCode, String)> {
    // 1. Compliance Check
    let risk_score = state.range.check_risk(&payload.destination_addr).await
        .map_err(|e| (StatusCode::INTERNAL_SERVER_ERROR, e))?;
    
    if risk_score > 80 && !payload.force.unwrap_or(false) {
        return Err((StatusCode::FORBIDDEN, "High risk destination address. Use --force to override (not recommended).".to_string()));
    }

    if risk_score > 80 && payload.force.unwrap_or(false) {
        println!("⚠️  WARNING: Bypassing Range Protocol firewall for high-risk address: {}", payload.destination_addr);
    }

    // 2. Route Selection
    let provider = if payload.strategy.contains("p2p") {
        &state.providers[1] // Radr
    } else {
        &state.providers[0] // Privacy Cash
    };

    // 3. Persist Intent
    let task_id = {
        let db = state.db.lock().unwrap();
        db.create_transaction(payload.amount_lamports, &payload.destination_addr, &provider.name())
            .map_err(|e| (StatusCode::INTERNAL_SERVER_ERROR, e.to_string()))?
    };

    // 4. Execution
    let result = provider.shield(payload, state.keystore.clone(), state.rpc.clone()).await
        .map_err(|e| {
            let db = state.db.lock().unwrap();
            let _ = db.update_status(&task_id, "Failed", None);
            (StatusCode::INTERNAL_SERVER_ERROR, e)
        })?;

    // 5. Update Status
    {
        let db = state.db.lock().unwrap();
        db.update_status(&task_id, "Confirmed", Some(&result.tx_hash))
            .map_err(|e| (StatusCode::INTERNAL_SERVER_ERROR, e.to_string()))?;
        
        if let Some(ref note) = result.note {
            db.update_note(&task_id, note)
                .map_err(|e| (StatusCode::INTERNAL_SERVER_ERROR, e.to_string()))?;
        }
    }

    Ok(Json(result))
}

pub async fn swap_handler(
    State(state): State<Arc<AppState>>,
    Json(payload): Json<SwapRequest>,
) -> Result<Json<SwapResponse>, (StatusCode, String)> {
    let result = state.swap_provider.swap(payload, state.keystore.clone(), state.rpc.clone()).await
        .map_err(|e| (StatusCode::INTERNAL_SERVER_ERROR, e))?;
    
    Ok(Json(result))
}

pub async fn pay_handler(
    State(state): State<Arc<AppState>>,
    Json(payload): Json<PayRequest>,
) -> Result<Json<PayResponse>, (StatusCode, String)> {
    let result = state.pay_provider.pay(payload, state.keystore.clone(), state.rpc.clone()).await
        .map_err(|e| (StatusCode::INTERNAL_SERVER_ERROR, e))?;
    
    Ok(Json(result))
}

pub async fn market_handler(
    State(state): State<Arc<AppState>>,
) -> Json<serde_json::Value> {
    let price = state.market.get_sol_price().await;
    Json(json!({
        "asset": "SOL",
        "price_usd": price,
        "provider": "Encrypt.trade"
    }))
}

pub async fn get_task_handler(

    State(state): State<Arc<AppState>>,

    Path(id): Path<String>,

) -> Result<Json<TransactionRecord>, (StatusCode, String)> {

    let db = state.db.lock().unwrap();

    let record = db.get_transaction(&id)

        .map_err(|_| (StatusCode::NOT_FOUND, "Task not found".to_string()))?;

    

    Ok(Json(record))

}



pub async fn get_history_handler(

    State(state): State<Arc<AppState>>,

) -> Result<Json<Vec<TransactionRecord>>, (StatusCode, String)> {

    let db = state.db.lock().unwrap();

    let records = db.list_transactions()

        .map_err(|e| (StatusCode::INTERNAL_SERVER_ERROR, e.to_string()))?;

    

    Ok(Json(records))

}
