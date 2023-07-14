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

static INSERT_POST_QUERY: &str = "INSERT INTO blog.posts (title, published_at, draft, content, path)
VALUES ($1, $2, $3, $4, $5) ON CONFLICT (path) DO UPDATE SET
title = EXCLUDED.title, published_at = EXCLUDED.published_at, draft = EXCLUDED.draft, content = EXCLUDED.content
RETURNING post_id";
static INSERT_TAGS_QUERY: &str = "INSERT INTO blog.tags (TAG) VALUES ($1) ON CONFLICT DO NOTHING";
static INSERT_POST_TAGS_QUERY: &str = "INSERT INTO blog.post_tags (post_id, tag) VALUES ($1, $2) ON CONFLICT DO NOTHING";
static REMOVE_UNUSED_POST_TAGS_QUERY: &str = "WITH unused AS (
SELECT tag FROM blog.tags NATURAL LEFT OUTER JOIN blog.post_tags WHERE post_id IS NULL
)
DELETE FROM blog.tags WHERE tag IN (SELECT tag from unused)";


// TODO(javier): Move function to lib.rs
fn store_post(client: &mut Client, post: blog::Post) {
    let stmt = client
        .prepare(INSERT_POST_QUERY)
        .expect("Could not prepare the statement");
    let rows = client
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
    let post_id: i32 = rows.iter().next().unwrap().get(0);
    let stmt = client.prepare(INSERT_TAGS_QUERY).expect("Could not prepare second statement");
    let insert_post_tags_stmt = client.prepare(INSERT_POST_TAGS_QUERY).unwrap();
    for tag in post.tags {
        let _ret = client.query(&stmt, &[&tag]).expect("Could not insert tag");
        let _ret = client.query(&insert_post_tags_stmt, &[&post_id, &tag]).unwrap();

    }

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
    let stmt = client.prepare(REMOVE_UNUSED_POST_TAGS_QUERY).unwrap();
    let _ret = client.query(&stmt, &[]).unwrap();
    Ok(())
}
