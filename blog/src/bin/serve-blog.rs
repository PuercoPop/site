//! # Serves the blog

use clap::Parser;

#[derive(Parser, Debug)]
struct Opts {
    #[arg(short = 'd')]
    dburl: String,
}

fn main() {
    let args = Opts::parse();
    let app = blog::new(args.dburl);

    tokio::runtime::Builder::new_multi_thread()
        .enable_all()
        .build()
        .expect("")
        .block_on(async {
            axum::Server::bind(&"0.0.0.0:3000".parse().unwrap())
                .serve(app.into_make_service())
                .await
                .unwrap();
        })
}
