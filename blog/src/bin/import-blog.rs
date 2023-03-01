// import-blog [directory]
use blog::read_post;
use clap::Parser;
use std::fs;
use std::io;
use std::path::Path;
use postgres::{Client, NoTls};

#[derive(Parser, Debug)]
struct Opts {
    // TODO(javier): Figure out how to pass values as environment variables as
    // well.
    #[arg(short = 'D')]
    dir: String,
    #[arg(short = 'd')]
    dburl: String,
}

fn main() -> Result<(), io::Error> {
    // 1. Iterate over the directory
    // 2. Filter out the posts that start with draft: in the header
    // 3. Insert the remaining posts using the path to determine whether the
    //    post already exists in the database.
    let args = Opts::parse();
    let dir = &args.dir.to_owned();
    let path = Path::new(dir);
    for entry in fs::read_dir(path)? {
        let entry = entry?;
        let metadata = entry.metadata()?;
        // TODO(javier): Replace with is_post()
        if !metadata.is_dir() {
            let post = read_post(&entry.path());
            println!("post: {post:#?}");
        }
    }

    let dburl = args.dburl;
    // TODO(javier): Enable TLS
    let mut _client = Client::connect(&dburl, NoTls).expect("Could not connect");
    Ok(())
}
