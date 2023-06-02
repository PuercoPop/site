use axum::{routing::get, Router, extract::State, response::Html};
use chrono::NaiveDate;
use minijinja::{context, Environment, Source};
use pulldown_cmark::{Event, Parser};
use regex::Regex;
use std::fs;
use std::io::{self, BufRead, BufReader};
use std::path::Path;
use std::sync::Arc;

// pub mod view;

#[derive(Debug, PartialEq)]
struct Tag {
    name: String,
}

#[derive(Debug, Default)]
pub struct Post {
    pub title: String,
    pub pubdate: NaiveDate,
    pub draft: bool,
    tags: Vec<Tag>,
    pub path: String,
}
impl Post {
    pub(crate) fn new() -> Post {
        Post::default()
    }
}

// Error return when read_post fails, explaining why it wasn't possible to read the post.
#[derive(Debug)]
pub enum PostParseError {
    IO(io::Error),
    CHRONO(chrono::ParseError),
    BadFormat,
}

impl From<io::Error> for PostParseError {
    fn from(err: io::Error) -> PostParseError {
        PostParseError::IO(err)
    }
}
impl From<chrono::ParseError> for PostParseError {
    fn from(err: chrono::ParseError) -> PostParseError {
        PostParseError::CHRONO(err)
    }
}

enum MetadataParseState {
    TitleLine,
    Tags,
    DateLine,
    End, // End of front matter
}

type FSM = MetadataParseState;

fn read_title(mut post: Post, line: &str) -> Result<Post, PostParseError> {
    let re = Regex::new(r"^(?:\s+)*Draft: (.*)")
        .expect("I need to update the error type so I can use ? ");
    let parser = Parser::new(line);
    for ev in parser {
        if let Event::Text(text) = ev {
            // Check if the title starts with Draft
            if re.is_match(&text) {
                let cap = re.captures(&text).expect("I suck");
                println!("cap: {cap:#?}");
                post.title = cap[1].to_string();
                post.draft = true;
            } else {
                post.title = text.to_string();
                post.draft = false;
            }

            return Ok(post);
        }
    }
    Err(PostParseError::BadFormat)
}

fn read_tags<'a>(mut post: Post, line: &'a str) -> Result<Post, PostParseError> {
    let parser = Parser::new(line);
    let mut tags: Vec<Tag> = Vec::new();
    for ev in parser {
        if let Event::Text(text) = ev {
            for tag in text.split(',') {
                // trim
                tags.push(Tag {
                    name: tag.trim().to_string(),
                })
            }
            post.tags = tags;
            return Ok(post);
        }
    }
    Err(PostParseError::BadFormat)
}

fn read_pubdate<'a>(mut post: Post, line: &'a str) -> Result<Post, PostParseError> {
    let parser = Parser::new(line);
    for ev in parser {
        if let Event::Text(text) = ev {
            let pubdate = NaiveDate::parse_from_str(&text, "%Y-%m-%d")?;
            post.pubdate = pubdate;
            return Ok(post);
        }
    }
    Err(PostParseError::BadFormat)
}

// Reads the meta-data embedded in the markdown document and returns a Post.
pub fn read_post(path: &Path) -> Result<Post, PostParseError> {
    let fd = fs::File::open(path)?;
    let reader = BufReader::new(fd);
    let mut post = Post::new();
    post.path = path.to_string_lossy().into_owned();
    let mut state: FSM = FSM::TitleLine;
    // TODO(javier): Handle EOF?
    for line in reader.lines() {
        let l = line?;
        match state {
            FSM::TitleLine => {
                post = read_title(post, &l)?;
                state = FSM::Tags;
            }
            FSM::Tags => {
                post = read_tags(post, &l)?;
                state = FSM::DateLine;
            }
            FSM::DateLine => {
                post = read_pubdate(post, &l)?;
                state = FSM::End;
            }
            FSM::End => return Ok(post),
        }
    }
    Err(PostParseError::BadFormat)
}

struct AppState {
    /// The template engine
    templates: Environment<'static>,
}

/// Initializes the application. Takes the URL of the database to use.
pub fn new(_dburl: String) -> Router {
    let source = Source::from_path("./templates");
    let mut env = Environment::new();
    env.set_source(source);
    let app_state = Arc::new(AppState { templates: env });
    let app = Router::new().route("/", get(index)).with_state(app_state);
    app
}

async fn index(State(state): State<Arc<AppState>>) -> Html<String> {
    let tmpl = state
        .templates
        .get_template("index.html")
        .expect("Unable to get template");
    Html(tmpl.render(context!()).expect("Unable to render template"))
}

#[cfg(test)]
mod tests {
    use super::*;
    #[test]
    fn test_read_title_1() {
        let post = Post::new();
        let line = "# Some title";
        let got = read_title(post, line).unwrap();
        assert_eq!(got.title, "Some title");
        assert_eq!(got.draft, false);
    }
    #[test]
    fn test_read_title_2() {
        let post = Post::new();
        let line = "# Draft: Some title";
        let got = read_title(post, line).unwrap();
        assert_eq!(got.title, "Some title");
        assert_eq!(got.draft, true);
    }
    #[test]
    fn test_read_tags_1() {
        let post = Post::new();
        let line = "## en, Emacs, rant";
        let got = read_tags(post, line).unwrap();
        let want: Vec<Tag> = vec![
            Tag {
                name: "en".to_string(),
            },
            Tag {
                name: "Emacs".to_string(),
            },
            Tag {
                name: "rant".to_string(),
            },
        ];
        assert_eq!(got.tags, want);
    }
    #[test]
    fn test_read_tags_2() {
        let post = Post::new();
        let line = "## en";
        let got = read_tags(post, line).unwrap();
        let want: Vec<Tag> = vec![Tag {
            name: "en".to_string(),
        }];
        assert_eq!(got.tags, want);
    }
    #[test]
    fn test_read_pubdate_1() {
        let post = Post::new();
        let line = "## 2022-02-15";
        let got = read_pubdate(post, line).unwrap();
        let want = NaiveDate::from_ymd_opt(2022, 2, 15).unwrap();
        assert_eq!(got.pubdate, want);
    }
    #[test]
    fn test_read_pubdate_2() {
        let post = Post::new();
        let line = "## 2022-2-15"; // w/o leading 0
        let got = read_pubdate(post, line).unwrap();
        let want = NaiveDate::from_ymd_opt(2022, 2, 15).unwrap();
        assert_eq!(got.pubdate, want);
    }
    #[test]
    fn test_read_pubdate_3() {
        let post = Post::new();
        let line = "## 2022-2-31"; // impossible date
        let got = read_pubdate(post, line);
        assert!(got.is_err());
    }

    #[test]
    fn test_read_post_0() {
        let path = Path::new("./testdata/post.md");
        let post = read_post(path).expect("Could not read post");
        assert_eq!(post.title, "Some title");
        assert_eq!(post.draft, false);
        assert_eq!(post.pubdate, NaiveDate::from_ymd_opt(2022, 3, 30).unwrap());
        assert_eq!(
            post.tags,
            vec![
                Tag {
                    name: "en".to_string(),
                },
                Tag {
                    name: "testing".to_string(),
                },
            ]
        );
        assert_eq!(post.path, path.to_str().unwrap().to_string());
    }

    #[test]
    #[ignore]
    fn test_read_post_1() {
        let path = Path::new("./testdata/draft_01.md");
        let post = read_post(path).expect("Could not read post");
        assert_eq!(post.draft, true)
    }
    #[test]
    #[ignore]
    fn test_read_post_2() {
        let path = Path::new("./testdata/post_01.md");
        let post = read_post(path).expect("Could not read post");
        assert_eq!(post.draft, false)
    }
}
