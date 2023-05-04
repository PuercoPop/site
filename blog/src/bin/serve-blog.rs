//! # Serves the blog

use blog::App;
use clap::Parser;

#[derive(Parser, Debug)]
struct Opts {
    #[arg(short = 'd')]
    dburl: String,
}

fn main() {
    let args = Opts::parse();
    let app = App::new(args.dburl);

    axum::Server::bind(&"0.0.0.0:3000".parse().unwrap())
        .serve(app.into_make_service())
        .await
        .unwrap();
}
