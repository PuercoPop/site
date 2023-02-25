// import-blog [directory]
use std::fs;
use clap::Parser;

#[derive(Parser, Debug)]
struct Opts {
     // TODO(javier): Figure out how to pass values as environment variables as
     // well.
    #[arg(short = "D")]
    dir: String,
    #[arg(short = "d")]
    dburl: String,

}

fn main() -> io::Result<()> {
    // 1. Iterate over the directory
    // 2. Filter out the posts that start with draft: in the header
    // 3. Insert the remaining posts using the path to determine whether the
    //    post already exists in the database.
    let args = Opts::parse();
    fs::read_dir(args.dir)?;

    // let (client, conn) = tokio_postgres::connect("", ).await?;

}
