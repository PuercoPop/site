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

#[tokio::test]
async fn test_authenticate_user_1() {
    // When no user is found
    let db = setup_db().await;
    let ret = authenticate_user(&db, "jane@doe.com".to_string(), "t0ps3cr3t".to_string())
        .await
        .unwrap();
    assert_eq!(ret, false);
}

#[tokio::test]
async fn test_authenticate_user_2() {
    // When the wrong path is provided
    let mut db = setup_db().await;
    let tx = db.transaction().await.unwrap();
    let count = tx
        .execute(
            "INSERT INTO finsta.users (email, password) VALUES ($1::text, crypt($2::text, gen_salt('bf', 8)))",
            &[&"jane@doe.com".to_string(), &"badtzu".to_string()],
        )
        .await.unwrap();
    assert_eq!(count, 1);
    let ret = authenticate_user(&tx, "jane@doe.com".to_string(), "t0ps3cr3t".to_string())
        .await
        .unwrap();
    assert_eq!(ret, false);
    tx.rollback().await.unwrap();
}

#[tokio::test]
async fn test_authenticate_user_3() {
    // Happy path
    let mut db = setup_db().await;
    let tx = db.transaction().await.unwrap();
    let count = tx
        .execute(
            "INSERT INTO finsta.users (email, password) VALUES ($1::text, crypt($2::text, gen_salt('bf', 8)))",
            &[&"jane@doe.com".to_string(), &"t0ps3cr3t".to_string()],
        )
        .await
        .unwrap();
    assert_eq!(count, 1);
    let ret = authenticate_user(&tx, "jane@doe.com".to_string(), "t0ps3cr3t".to_string())
        .await
        .unwrap();
    assert_eq!(ret, true);
    tx.rollback().await.unwrap()
}
