use clap::Parser;
use finsta::new;

#[derive(Parser)]
struct Opts {
    #[arg(short = 'D')]
    templates_dir: String,
}

#[tokio::main]
async fn main() {
    let args = Opts::parse();
    let app = new(args.templates_dir);
    axum::Server::bind(&"0.0.0.0:3000".parse().unwrap())
        .serve(app.into_make_service())
        .await
        .unwrap();
}
