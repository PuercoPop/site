use tokio_postgres::{Client, Error as PgError};

static CHECK_PASS: &str =
    "SELECT u.password = crypt($1::text, u.password) AS result FROM finsta.users u WHERE email = $2";

pub(crate) async fn authenticate_user(
    db: &Client,
    email: String,
    password: String,
) -> Result<bool, PgError> {
    let stmt = db.prepare(CHECK_PASS).await?;
    let row = db.query_one(&stmt, &[&password, &email]).await?;
    let ret: bool = row.get::<&str, bool>("result");
    Ok(ret)
}
