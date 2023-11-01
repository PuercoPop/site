use clap::Parser;
use finsta::new;

#[derive(Parser)]
struct Opts {
    #[arg(short = 'D')]
    templates_dir: String,
}

#[derive(thiserror::Error, Debug)]
enum Error {
    #[error(transparent)]
    HyperError(#[from] hyper::Error)
}

#[tokio::main]
async fn main() -> Result<(), Error> {
    let args = Opts::parse();
    let app = new(args.templates_dir);
    axum::Server::bind(&"0.0.0.0:3000".parse().unwrap())
        .serve(app.into_make_service())
        .await?;
    Ok(())
}
