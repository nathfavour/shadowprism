pub mod db;
pub mod keystore;
pub mod middleware;
pub mod adapters;
pub mod api;

use axum::{
    routing::{get, post},
    Router,
};
use std::sync::Arc;
use crate::api::{AppState, shield_handler};
use crate::adapters::{privacy_cash::PrivacyCashAdapter, radr::RadrAdapter, range::RangeClient};

#[cfg(unix)]
use tokio::net::UnixListener;

#[tokio::main]
async fn main() {
    tracing_subscriber::fmt::init();

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

    // Resolve ~/.shadowprism/engine.sock
    let mut uds_path = home::home_dir().expect("Could not find home directory");
    uds_path.push(".shadowprism");
    
    // Ensure directory exists
    if !uds_path.exists() {
        std::fs::create_dir_all(&uds_path).expect("Could not create .shadowprism directory");
    }
    
    uds_path.push("engine.sock");

    #[cfg(unix)]
    {
        let _ = std::fs::remove_file(&uds_path);
        let listener = UnixListener::bind(&uds_path).expect("Could not bind to UDS socket");
        
        // Set permissions so only the user can access the socket
        use std::os::unix::fs::PermissionsExt;
        std::fs::set_permissions(&uds_path, std::fs::Permissions::from_mode(0o700)).unwrap();
        
        println!("🛡️ ShadowPrism Core online via UDS: {:?}", uds_path);
        axum::serve(listener, app).await.unwrap();
    }

    #[cfg(not(unix))]
    {
        let port = std::env::var("PORT").unwrap_or_else(|_| "42069".to_string());
        let addr = format!("127.0.0.1:{}", port).parse::<std::net::SocketAddr>().unwrap();
        let listener = tokio::net::TcpListener::bind(&addr).await.unwrap();
        println!("🛡️ ShadowPrism Core online via TCP: {}", addr);
        axum::serve(listener, app).await.unwrap();
    }
}

async fn health_check() -> &'static str {
    "{\"status\": \"ready\", \"engine\": \"rust\", \"protocol\": \"uds\"}"
}