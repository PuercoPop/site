use axum::{
    routing::{get, post},
    Router,
};
use minijinja::{path_loader, Environment};
use std::sync::Arc;
use tokio_postgres::Client;

mod db;
mod users;
mod handlers;

/// A context that is accessible on all axum request handlers.
pub struct HTTPContext {
    // The connection to the database
    db: Client,
    /// An environment containing the HTML templates
    templates: Environment<'static>,
}

pub fn new(template_dir: String, db: Client) -> Router {
    let mut env = Environment::new();
    env.set_loader(path_loader(template_dir));
    let ctx = Arc::new(HTTPContext {
        templates: env,
        db: db,
    });

    Router::new()
        .route("/", get(handlers::index))
        .route("/sign-in", post(handlers::sign_in))
        .route("/upload", get(handlers::form))
        .route("/upload", post(handlers::upload))
        .with_state(ctx)
}
