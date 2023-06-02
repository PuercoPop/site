//! # Serves the blog

use clap::Parser;

#[derive(Parser, Debug)]
struct Opts {
    #[arg(short = 'd')]
    dburl: String,
    #[arg(short = 'D')]
    templates_dir: String
}

fn main() {
    let args = Opts::parse();
    let ctx = blog::new_ctx(args.dburl, args.templates)?
    let app = blog::new();

    tokio::runtime::Builder::new_multi_thread()
        .enable_all()
        .build()
        .expect("[serve-blog]: Unable to build the tokio runtime")
        .block_on(async {
            axum::Server::bind(&"0.0.0.0:3000".parse().unwrap())
                .serve(app.into_make_service())
                .await
                .unwrap();
        })
}
