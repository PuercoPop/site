// -*- compile-command: rustc migrate-rs.rs-*-
//! # Design
//! 1. List all the migrations in directory
//! 2. Run a psql query to find out which ones have not been run.
//! 3. For each migration
//!    a. Run psql -f <path/to/migration>
//!    b. Check the status code
//!      i.  If 0, record the migration
//!      ii. If non-zero abort the entire process
//!
//! How should we communicate when the hash of the migration has changed?
//! Use PostgreSQL JSON support to output structured data. We can use jq to process
//! it.
//!
//! How do we test the CLI util?
//! - expect?
//! - go?
//!
//! What are the scenarios to test?
use std::env;
use std::process::exit;

fn print_help() {
    println!(r#"Usage: migrate-rs [-h] -D <migrations dir> -d <database URL>"#);
}

fn main() {
    let mut migration_dir: Option<String> = None;
    let mut dburl: Option<String> = None;
    let mut args = env::args().into_iter();
    // drop the first one argument which is the name of the executable
    args.next();
    while let Some(arg) = args.next() {
        match arg.as_str() {
            "-h" => {
                print_help();
                exit(0)
            }
            "-D" => match args.next() {
                Some(v) => migration_dir = Some(v),
                None => {
                    println!("-D takes a directory");
                    print_help();
                    exit(1);
                }
            },
            "-d" => {
                match args.next() {
                    Some(v) => dburl = Some(v),
                    None => {
                        println!("-D takes the URL of the PostgreSQL database to connect to");
                        print_help();
                        exit(1);
                    }
                }
                println!("database URL to connect to")
            }
            &_ => {
                println!("Unrecognized option: {:?}", arg);
                print_help();
                exit(1);
            }
        }
    }
}
