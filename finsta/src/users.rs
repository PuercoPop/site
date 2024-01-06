use tokio_postgres::{Error as PgError, GenericClient};

static CHECK_PASS: &str =
    "SELECT u.password = crypt($1::text, u.password) AS result FROM finsta.users u WHERE email = $2";

pub async fn authenticate_user<C: GenericClient>(
    db: &C,
    email: String,
    password: String,
) -> Result<bool, PgError> {
    let stmt = db.prepare(CHECK_PASS).await?;
    // TODO: Can I use query_opt(…).and_then(…) instead?
    let rows = db.query(&stmt, &[&password, &email]).await?;
    // If the rows <> 1 return false.
    if rows.len() != 1 {
        return Ok(false);
    }
    rows[0].try_get("result")
}
