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
use std::fs;
use std::io::Write;
use std::path::PathBuf;
use std::process::{exit, Command};

fn print_help() {
    println!(r#"Usage: migrate-rs [-h] -D <migrations dir> -d <database URL>"#);
}

fn main() -> Result<(), Error> {
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
    let dburl = dburl.unwrap();
    let migration_dir = migration_dir.unwrap();
    let files = list_migrations(migration_dir)?;
    ensure_migration_schema(&dburl)?;
    for m in files {
        let res = run_migration(&dburl, &m);
        match res {
            Ok(success) => {
                exit(if success { 0 } else { 1 });
            }
            Err(_) => {
                exit(1);
            }
        }
    }
    Ok(())
}

#[derive(Debug)]
enum Error {
    /// Catch all value for any errors
    Any,
    /// Signals that an error occurred when installing the migration schema.
    SchemaFailure,
    /// An error was encountered while trying to execute a migration.
    MigrationFailure,
}

/// Returns an ordered list of paths to SQL files.
fn list_migrations(dir: String) -> Result<Vec<PathBuf>, Error> {
    // TODO: Filter non .sql files
    let files = fs::read_dir(dir)
        .or(Err(Error::Any))
        .map(|rdir| {
            rdir.map(|direntry| match direntry {
                Ok(v) => Ok(v.path()),
                Err(_) => Err(Error::Any),
            })
        })?
        .collect();
    files
}

/// Ensures that the migration schema is present in the database at `dburl`.
fn ensure_migration_schema(dburl: &String) -> Result<(), Error> {
    let output = Command::new("psql")
        .arg("-d")
        .arg(dburl)
        .arg("-f")
        .arg("meta.sql")
        .output()
        .or(Err(Error::Any))?;
    if output.status.success() {
        Ok(())
    } else {
        Err(Error::SchemaFailure)
    }
}

/// Execute migration at `path`. Abort program if migration couldn't be applied successfully.
fn run_migration(dburl: &String, path: &PathBuf) -> Result<bool, Error> {
    let output = Command::new("sql")
        .arg("-d")
        .arg(dburl)
        .arg("f")
        .arg(path)
        .output()
        .or(Err(Error::MigrationFailure))?;
    if output.status.success() {
        // Call record migration here
        record_migration(dburl, &path).or(Err(Error::MigrationFailure))
    } else {
        Err(Error::MigrationFailure)
    }
}

enum MigrationError {
    /// An error occurred while trying to compute the checksum of the migration
    Checksum,
    /// Could not convert the path to a string
    Filename,
    /// Could not spawn psql
    Spawn,
}

/// Record the fact a migration was successfully run.
fn record_migration(dburl: &String, path: &PathBuf) -> Result<bool, MigrationError> {
    // string builder
    let checksum = md5(path).ok_or(MigrationError::Checksum)?;
    let path = path.to_str().ok_or(MigrationError::Filename)?;
    let mut sql = String::new();
    sql.push_str("INSERT INTO public.versions (version, checksum) VALUES (");
    sql.push_str(path);
    sql.push_str(", ");
    sql.push_str(checksum.as_str());
    sql.push_str(" )");
    // psql put stdin
    // https://doc.rust-lang.org/std/process/struct.Stdio.html#method.piped
    let mut child = Command::new("psql")
        .arg("-d")
        .arg(dburl)
        .spawn()
        .or(Err(MigrationError::Spawn))?;
    let mut stdin = child.stdin.take().expect("Failed to open stdin");
    stdin.write_all(sql.as_bytes()).expect("Foo");
    drop(stdin);
    let code = child.wait().expect("Could not wait");
    Ok(code.success())
}

/// Returns the md5 checksum of the file at `path`.
fn md5(path: &PathBuf) -> Option<String> {
    path.to_str().and_then(|path| {
        Command::new("md5sum")
            .arg(path)
            .output()
            .ok()
            .filter(|out| out.status.success())
            .map(|out| out.stdout)
            .and_then(|bytes| String::from_utf8(bytes).ok())
            .and_then(|cmdout| cmdout.split(' ').next().map(|out| out.to_string()))
    })
}

#[cfg(test)]
mod tests {
    use super::*;

    // #[test]
    // fn test_record_migration_sql() {
    //     let pathbuf = PathBuf::new("./meta.sql");
    //     assert_eq!(
    //         test_record_migration_sql(pathbuf),
    //         "INSERT INTO public.versions (version, checksum) VALUES (meta.sql, 'foo')"
    //     );
    // }

    #[test]
    fn test_md5_1() {
        let pathbuf = PathBuf::from("./meta.sql");
        assert_eq!(
            md5(&pathbuf),
            Some("a9c1d60ea6e7d398c30f31c854a1a62c".to_string())
        );
    }

    #[test]
    fn test_md5_enoent() {
        let pathbuf = PathBuf::from("./fubar.sql");
        assert_eq!(md5(&pathbuf), None);
    }
}
