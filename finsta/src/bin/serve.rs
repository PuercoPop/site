use clap::Parser;
use finsta::new;
use tokio_postgres::NoTls;
use tower_http::trace::TraceLayer;

#[derive(Parser)]
struct Opts {
    #[arg(short = 'd')]
    dburl: String,
    #[arg(short = 'D')]
    templates_dir: String,
}

#[derive(thiserror::Error, Debug)]
enum Error {
    #[error(transparent)]
    HyperError(#[from] hyper::Error),
    #[error(transparent)]
    DBError(#[from] tokio_postgres::Error),
}

#[tokio::main]
async fn main() -> Result<(), Error> {
    let args = Opts::parse();
    let dburl = args.dburl.as_str();
    let (client, conn) = tokio_postgres::connect(dburl, NoTls).await?;
    tokio::spawn(async move {
        if let Err(err) = conn.await {
            eprintln!("connection error: {}", err);
        }
    });
    tracing_subscriber::fmt::init();

    let app = new(args.templates_dir, client)
        .layer(TraceLayer::new_for_http());
    axum::Server::bind(&"0.0.0.0:3000".parse().unwrap())
        .serve(app.into_make_service())
        .await?;
    Ok(())
}
