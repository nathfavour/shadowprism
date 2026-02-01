pub mod db;
pub mod keystore;
pub mod middleware;
pub mod adapters;
pub mod api;

use axum::{
    routing::{get, post},
    Router,
};
use std::net::SocketAddr;
use std::sync::Arc;
use crate::api::{AppState, shield_handler};
use crate::adapters::{privacy_cash::PrivacyCashAdapter, radr::RadrAdapter, range::RangeClient};

#[tokio::main]
async fn main() {
    tracing_subscriber::fmt::init();

    // Setup state
    let state = Arc::new(AppState {
        range: RangeClient::new(),
        providers: vec![
            Box::new(PrivacyCashAdapter),
            Box::new(RadrAdapter),
        ],
    });

    let app = Router::new()
        .route("/health", get(health_check))
        .route("/v1/shield", post(shield_handler))
        .layer(axum::middleware::from_fn(middleware::auth_validator))
        .with_state(state);

    let port = std::env::var("PORT")
        .unwrap_or_else(|_| "42069".to_string())
        .parse::<u16>()
        .unwrap();

    let addr = SocketAddr::from(([127, 0, 0, 1], port));
    println!("🛡️ ShadowPrism Core online at {}", addr);
    
    let listener = tokio::net::TcpListener::bind(&addr).await.unwrap();
    axum::serve(listener, app).await.unwrap();
}

async fn health_check() -> &'static str {
    "{\"status\": \"ready\", \"engine\": \"rust\"}"
}