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
use rand::Rng;
use crate::api::{AppState, shield_handler, get_task_handler, get_history_handler, swap_handler, pay_handler, market_handler};
use crate::adapters::{
    privacy_cash::PrivacyCashAdapter, 
    radr::RadrAdapter, 
    range::RangeClient,
    market::MarketOracle,
    silent_swap::SilentSwapAdapter,
    starpay::StarpayAdapter,
    rpc::ReliableClient,
};
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
    let db_store = Arc::new(Mutex::new(store));

    // Initialize Keystore
    let mut master_key_path = data_dir.clone();
    master_key_path.push(".master");

    let passphrase = if master_key_path.exists() {
        std::fs::read_to_string(&master_key_path).expect("Failed to read master key")
    } else {
        let new_key: String = rand::thread_rng()
            .sample_iter(&rand::distributions::Alphanumeric)
            .take(64)
            .map(char::from)
            .collect();
        std::fs::write(&master_key_path, &new_key).expect("Failed to initialize master key");
        
        // Set restrictive permissions on the master key file
        #[cfg(unix)]
        {
            use std::os::unix::fs::PermissionsExt;
            std::fs::set_permissions(&master_key_path, std::fs::Permissions::from_mode(0o600)).unwrap();
        }
        new_key
    };

    let mut keystore_path = data_dir.clone();
    keystore_path.push("wallet.prism");
    
    let keystore = Arc::new(crate::keystore::PrismKeystore::load_or_create(&keystore_path, &passphrase)
        .expect("Failed to initialize encrypted keystore"));

    let rpc = Arc::new(ReliableClient::new());

    let state = Arc::new(AppState {
        rpc,
        range: RangeClient::new(),
        market: MarketOracle::new(),
        providers: vec![
            Box::new(PrivacyCashAdapter),
            Box::new(RadrAdapter),
        ],
        swap_provider: Box::new(SilentSwapAdapter),
        pay_provider: Box::new(StarpayAdapter),
        db: db_store.clone(),
        keystore,
    });

    // Start Watchdog
    let watchdog_store = db_store.clone();
    tokio::spawn(async move {
        let watchdog = crate::db::watchdog::Watchdog::new(watchdog_store);
        watchdog.start().await;
    });

    let app = Router::new()
        .route("/health", get(health_check))
        .route("/v1/shield", post(shield_handler))
        .route("/v1/swap", post(swap_handler))
        .route("/v1/pay", post(pay_handler))
        .route("/v1/market", get(market_handler))
        .route("/v1/tasks/{id}", get(get_task_handler))
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
