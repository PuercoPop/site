use finsta::new;

#[tokio::main]
async fn main() {
    // TODO: Extract config parameters as args/enviroment variables
    let app = new();
    axum::Server::bind(&"0.0.0.0:3000".parse().unwrap())
        .serve(app.into_make_service())
        .await
        .unwrap();
}
