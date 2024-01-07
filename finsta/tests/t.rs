use tokio_postgres::Client;

pub async fn new_db() -> Client {
    let dburl = std::env::var("FINSTA_TEST_DB")
        .unwrap_or("postgres://postgres@localhost:5432/finsta_test".to_string());

    let (client, conn) = tokio_postgres::connect(dburl.as_str(), tokio_postgres::NoTls)
        .await
        .expect("Could not connect to test database");
    tokio::spawn(async move {
        if let Err(err) = conn.await {
            eprintln!("Could not connect to test database: {}", err);
        }
    });
    client
}

// TODO: Implement with_rollback(|tx| ...)
