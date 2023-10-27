use axum::{routing::{get, post}, Router};
use finsta::{form, index, upload};

#[tokio::main]
async fn main() {
    // TODO: Extract config parameters as args/enviroment variables
    let app = Router::new()
        .route("/", get(index))
        .route("/upload", get(form))
        .route("/upload", post(upload));
    axum::Server::bind(&"0.0.0.0:3000".parse().unwrap())
        .serve(app.into_make_service())
        .await
        .unwrap();
}
