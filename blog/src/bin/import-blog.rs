// import-blog [directory]
use blog::{read_post, remove_unused_tags, store_post, PostParseError};
use clap::Parser;
use std::fs;
use std::io;
use std::path::Path;
use tokio_postgres::NoTls;

#[derive(Parser, Debug)]
struct Opts {
    // TODO(javier): Figure out how to pass values as environment variables as
    // well.
    #[arg(short = 'D')]
    dir: String,
    #[arg(short = 'd')]
    dburl: String,
}

#[derive(thiserror::Error, Debug)]
enum Error {
    #[error(transparent)]
    IoError(#[from] io::Error),
    #[error(transparent)]
    InvalidPost(#[from] PostParseError),
    #[error(transparent)]
    PgError(#[from] tokio_postgres::Error),
}

#[tokio::main]
async fn main() -> Result<(), Error> {
    // 1. Iterate over the directory
    // 2. Filter out the posts that start with draft: in the header
    // 3. Insert the remaining posts using the path to determine whether the
    //    post already exists in the database.
    let args = Opts::parse();
    let dburl = args.dburl;
    // TODO(javier): Enable TLS
    let (client, conn) = tokio_postgres::connect(&dburl, NoTls).await?;
    let dir = &args.dir.to_owned();
    let path = Path::new(dir);

    tokio::spawn(async move {
        if let Err(err) = conn.await {
            eprintln!("connection error: {}", err);
        }
    });

    // TODO: The code below should be grouped in a single import_posts(dir:
    // Path) -> Result<(), Error> call. Make store_post, remove_unused_tags.
    for entry in fs::read_dir(path)? {
        let entry = entry?;
        let metadata = entry.metadata()?;
        // TODO(javier): Replace with is_post()
        if !metadata.is_dir() {
            let post = read_post(&entry.path())?;
            if !post.draft {
                store_post(&client, post).await?;
            }
        }
    }
    remove_unused_tags(&client).await?;

    Ok(())
}
