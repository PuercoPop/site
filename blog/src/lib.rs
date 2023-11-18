use axum::{routing::get, Router};
use minijinja::{Environment, path_loader};
use std::sync::Arc;
use tokio_postgres::Client;

mod db;
mod handlers;
mod post;

pub use self::db::{remove_unused_tags, store_post};
pub use self::post::{read_post, PostParseError};

/// Holds to the global resources the application depends on
// TODO: Rename to WebContext
pub struct Context {
    /// The template engine
    templates: Environment<'static>,
    /// A connection to the database
    db: Client,
}

pub fn new_ctx(client: Client, template_dir: String) -> Context {
    let source = path_loader(template_dir);
    let mut env = Environment::new();
    env.set_loader(source);

    Context {
        templates: env,
        db: client,
    }
}

/// Initializes the application. Takes the URL of the database to use.
pub fn new(ctx: Context) -> Router {
    // TODO: Can I remove the Arc wrapper?
    let ctx = Arc::new(ctx);
    Router::new()
        .route("/", get(handlers::index))
        .route("/p/:slug", get(handlers::show_post))
        .route("/tags/", get(handlers::list_tags))
        .route("/t/:tag", get(handlers::show_tag))
        .route("/feed/", get(handlers::feed))
        .with_state(ctx)
}
