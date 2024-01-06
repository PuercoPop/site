use tokio_postgres::{Client, Error as PgError};

static CHECK_PASS: &str =
    "SELECT u.password = crypt($1::text, u.password) AS result FROM finsta.users u WHERE email = $2";

pub async fn authenticate_user(
    db: &Client,
    email: String,
    password: String,
) -> Result<bool, PgError> {
    let stmt = db.prepare(CHECK_PASS).await?;
    let rows = db.query(&stmt, &[&password, &email]).await?;
    // If the rows <> 1 return false.
    if rows.len() != 1 {
        return Ok(false);
    }
    rows[0].try_get("result")
}
