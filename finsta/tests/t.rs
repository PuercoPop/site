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

// TODO: Finish this method once async closures become stable. Ref:
// https://github.com/rust-lang/rust/issues/62290
//
// pub async fn with_rollback<F>(db: &mut Client, block: F)
// where
//     F: Fn(&Transaction) -> (),
// {
//     let tx = db.transaction().await.unwrap();
//     block(&tx);
//     tx.rollback().await.unwrap();
// }
