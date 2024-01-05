use finsta::users::authenticate_user;
use tokio_postgres::Client;

async fn setup_db() -> Client {
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

// TODO: Test authenticate_user with no users
#[tokio::test]
async fn test_authenticate_user_1() {
    let db = setup_db().await;
    let ret = authenticate_user(&db, "jane@doe.com".to_string(), "t0ps3cr3t".to_string())
        .await
        .unwrap();
    assert_eq!(ret, false);
}

// TODO: Test authenticate_user with wrong password
// TODO: Test authenticate_user with correct password
