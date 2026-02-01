pub mod db;
pub mod keystore;
pub mod middleware;
pub mod adapters;
pub mod api;

use axum::{
    routing::{get, post},
    Router,
};
use std::sync::{Arc, Mutex};
use crate::api::{AppState, shield_handler, get_task_handler, get_history_handler};
use crate::adapters::{privacy_cash::PrivacyCashAdapter, radr::RadrAdapter, range::RangeClient};
use crate::db::TransactionStore;

#[cfg(unix)]
use tokio::net::UnixListener;

#[tokio::main]
async fn main() {
    tracing_subscriber::fmt::init();

    // Resolve ~/.shadowprism/
    let mut data_dir = home::home_dir().expect("Could not find home directory");
    data_dir.push(".shadowprism");
    if !data_dir.exists() {
        std::fs::create_dir_all(&data_dir).expect("Could not create .shadowprism directory");
    }

    // Initialize Database
    let mut db_path = data_dir.clone();
    db_path.push("prism.db");
    let store = TransactionStore::new(&db_path).expect("Failed to initialize SQLite database");

    let state = Arc::new(AppState {
        range: RangeClient::new(),
        providers: vec![
            Box::new(PrivacyCashAdapter),
            Box::new(RadrAdapter),
        ],
        db: Mutex::new(store),
    });

    let app = Router::new()
        .route("/health", get(health_check))
        .route("/v1/shield", post(shield_handler))
        .route("/v1/tasks/:id", get(get_task_handler))
        .route("/v1/history", get(get_history_handler))
        .layer(axum::middleware::from_fn(middleware::auth_validator))
        .with_state(state);

    let mut uds_path = data_dir.clone();
    uds_path.push("engine.sock");

    #[cfg(unix)]
    {
        let _ = std::fs::remove_file(&uds_path);
        let listener = UnixListener::bind(&uds_path).expect("Could not bind to UDS socket");
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
