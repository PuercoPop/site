// import-blog [directory]
use blog::{read_post, PostParseError};
use clap::Parser;
use postgres::{Client, NoTls};
use std::fs;
use std::io;
use std::path::Path;

#[derive(Parser, Debug)]
struct Opts {
    // TODO(javier): Figure out how to pass values as environment variables as
    // well.
    #[arg(short = 'D')]
    dir: String,
    #[arg(short = 'd')]
    dburl: String,
}

// TODO(javier): Where does this function blog to?
fn store_post(client: &mut Client, post: blog::Post) {
    // TODO(javier): Set ON CONFLICT DO UPDATE
    let stmt = client
        .prepare(
            "INSERT INTO blog.posts (title, published_at, draft, content, path)
      VALUES ($1, $2, $3, $4, $5)",
        )
        .expect("Could not prepare the statement");
    let _ret = client
        .query(
            &stmt,
            &[
                &post.title,
                &post.pubdate,
                &post.draft,
                &post.content,
                &post.path,
            ],
        )
        .expect("Unable to insert post entry.");
    // TODO: then insert each tag
}

#[derive(Debug)]
enum Error {
    IoError(io::Error),
    InvalidPost(PostParseError),
}

impl From<std::io::Error> for Error {
    fn from(err: std::io::Error) -> Self {
        Error::IoError(err)
    }
}

impl From<PostParseError> for Error {
    fn from(err: PostParseError) -> Self {
        Error::InvalidPost(err)
    }
}

fn main() -> Result<(), Error> {
    // 1. Iterate over the directory
    // 2. Filter out the posts that start with draft: in the header
    // 3. Insert the remaining posts using the path to determine whether the
    //    post already exists in the database.
    let args = Opts::parse();
    let dburl = args.dburl;
    // TODO(javier): Enable TLS
    let mut client = Client::connect(&dburl, NoTls).expect("Could not connect");
    let dir = &args.dir.to_owned();
    let path = Path::new(dir);
    for entry in fs::read_dir(path)? {
        let entry = entry?;
        let metadata = entry.metadata()?;
        // TODO(javier): Replace with is_post()
        if !metadata.is_dir() {
            let post = read_post(&entry.path())?;
            if !post.draft {
                store_post(&mut client, post);
            }
        }
    }

    Ok(())
}
