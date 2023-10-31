use axum::{
    routing::{get, post},
    Router,
};
use minijinja::{path_loader, Environment};
use std::sync::Arc;

mod db;
mod handlers;

/// A context that is accessible on all axum request handlers.
pub struct HTTPContext {
    /// An environment containing the HTML templates
    templates: Environment<'static>,
}

pub fn new(template_dir: String) -> Router {
    let mut env = Environment::new();
    env.set_loader(path_loader(template_dir));
    let ctx = Arc::new(HTTPContext { templates: env });

    Router::new()
        .route("/", get(handlers::index))
        .route("/upload", get(handlers::form))
        .route("/upload", post(handlers::upload))
        .with_state(ctx)
}
