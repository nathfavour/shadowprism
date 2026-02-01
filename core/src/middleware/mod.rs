use axum::{
    body::Body,
    http::{Request, StatusCode},
    middleware::Next,
    response::Response,
};
use std::env;

pub async fn auth_validator(
    req: Request<Body>,
    next: Next,
) -> Result<Response, StatusCode> {
    let auth_header = req.headers()
        .get("Authorization")
        .and_then(|h| h.to_str().ok());

    let expected_token = env::var("SHADOWPRISM_AUTH_TOKEN")
        .unwrap_or_else(|_| "dev-token-123".to_string());
    
    let expected_auth = format!("Bearer {}", expected_token);

    if let Some(auth) = auth_header {
        if auth == expected_auth {
            return Ok(next.run(req).await);
        }
    }

    Err(StatusCode::UNAUTHORIZED)
}