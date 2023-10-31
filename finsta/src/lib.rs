use axum::{
    routing::{get, post},
    Router,
};

mod db;
mod handlers;

pub fn new() -> Router {
    Router::new()
        .route("/", get(handlers::index))
        .route("/upload", get(handlers::form))
        .route("/upload", post(handlers::upload))
}
