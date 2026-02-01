use axum::{
    extract::{State, Path},
    http::StatusCode,
    Json,
};
use crate::adapters::{ShieldRequest, ShieldResponse};
use crate::db::{TransactionStore, TransactionRecord};
use std::sync::{Arc, Mutex};

pub struct AppState {
    pub range: crate::adapters::range::RangeClient,
    pub providers: Vec<Box<dyn crate::adapters::PrivacyProvider>>,
    pub db: Mutex<TransactionStore>,
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

    // 3. Persist Intent
    let task_id = {
        let db = state.db.lock().unwrap();
        db.create_transaction(payload.amount_lamports, &payload.destination_addr, &provider.name())
            .map_err(|e| (StatusCode::INTERNAL_SERVER_ERROR, e.to_string()))?
    };

    // 4. Execution
    let result = provider.shield(payload).await
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
    }

    Ok(Json(result))
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
