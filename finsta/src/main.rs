use axum::{response::Html, routing::get, Router};

async fn index() -> Html<&'static str> {
    Html("こんにちは")
}

#[tokio::main]
async fn main() {
    let app = Router::new().route("/", get(index));
    axum::Server::bind(&"0.0.0.0:3000".parse().unwrap())
        .serve(app.into_make_service())
        .await
        .unwrap();
}
