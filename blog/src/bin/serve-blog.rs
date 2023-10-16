//! # Serves the blog

use clap::Parser;
use std::fmt;
use tokio_postgres::NoTls;

#[derive(Parser, Debug)]
struct Opts {
    #[arg(short = 'd')]
    dburl: String,
    #[arg(short = 'D')]
    templates_dir: String,
}

/// Errors that can happen during the initialization of the application.
#[derive(thiserror::Error, Debug)]
enum Error {
    /// Returned when we were unable to connect to the database
    DBerror(#[from] tokio_postgres::Error),
    /// Returned when we were unable to start the axum server
    ServerError(#[from] hyper::Error),
}

impl fmt::Display for Error {
    fn fmt(&self, fmt: &mut fmt::Formatter<'_>) -> fmt::Result {
        let msg = match &*self {
            Error::DBerror(_err) => "DBerror",
            Error::ServerError(_err) => "ServerError",
        };
        fmt.write_str(msg)
    }
}

#[tokio::main]
async fn main() -> Result<(), Error> {
    let args = Opts::parse();
    let dburl: &str = args.dburl.as_str();
    let (client, conn) = tokio_postgres::connect(dburl, NoTls).await?;
    let ctx = blog::new_ctx(client, args.templates_dir);
    let app = blog::new(ctx);

    tokio::spawn(async move {
        if let Err(err) = conn.await {
            eprintln!("connection error: {}", err);
        }
    });

    // TODO: Take the Host and Port as arguments
    match axum::Server::bind(&"0.0.0.0:3000".parse().unwrap())
        .serve(app.into_make_service())
        .await
    {
        Ok(()) => (),
        Err(err) => eprintln!("connection error: {}", err),
    }
    Ok(())
}
