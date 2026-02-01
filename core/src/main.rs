pub mod db;
pub mod keystore;
pub mod middleware;
pub mod adapters;
pub mod api;

use axum::{routing::get, Router};
use std::net::SocketAddr;

#[tokio::main]
async fn main() {
    // Initialize tracing
    tracing_subscriber::fmt::init();

    // Build our application with a route
    let app = Router::new()
        .route("/health", get(health_check));

    // Run it with hyper
    let addr = SocketAddr::from(([127, 0, 0, 1], 42069));
    println!("ShadowPrism Core listening on {}", addr);
    
    axum::Server::bind(&addr)
        .serve(app.into_make_service())
        .await
        .unwrap();
}

async fn health_check() -> &'static str {
    "{\"status\": \"ready\", \"engine\": \"rust\"}"
}