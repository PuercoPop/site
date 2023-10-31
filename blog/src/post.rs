use chrono::NaiveDate;
use pulldown_cmark::{Event, Parser};
use regex::Regex;
use serde::Serialize;
use std::io::{self, BufRead, BufReader};
use std::path::Path;
use std::{fs, io::Read};

pub type Tag = String;

#[derive(Serialize, Debug, Default)]
pub struct Post {
    pub slug: String,
    pub title: String,
    pub pubdate: NaiveDate,
    pub draft: bool,
    pub tags: Vec<Tag>,
    // /// The markdown source
    // pub source: String,
    /// The rendered HTML, minutes the embedded metadata
    pub content: String,
    pub path: String,
}

impl Post {
    pub(crate) fn new() -> Post {
        Post::default()
    }
}

// Error return when read_post fails, explaining why it wasn't possible to read
// the post
#[derive(thiserror::Error, Debug)]
pub enum PostParseError {
    #[error(transparent)]
    IO(#[from] io::Error),
    #[error(transparent)]
    CHRONO(#[from] chrono::ParseError),
    #[error("bad post format")]
    BadFormat,
}

enum MetadataParseState {
    TitleLine,
    Tags,
    DateLine,
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
                tags.push(tag.trim().to_string())
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
    // TODO(javier): We need to update this method to switch from reading line
    // by line, to reading the full content.
    let fd = fs::File::open(path)?;
    let mut reader = BufReader::new(fd);
    let mut post = Post::new();
    post.path = path.to_string_lossy().into_owned();
    let mut state: FSM = FSM::TitleLine;
    // TODO(javier): Handle EOF?
    loop {
        let mut line = String::new();
        let _len = reader.read_line(&mut line)?;
        match state {
            FSM::TitleLine => {
                post = read_title(post, &line)?;
                state = FSM::Tags;
            }
            FSM::Tags => {
                post = read_tags(post, &line)?;
                state = FSM::DateLine;
            }
            FSM::DateLine => {
                post = read_pubdate(post, &line)?;
                break;
            }
        }
    }
    let mut content = String::new();
    let mut body = String::new();
    let _len = reader.read_to_string(&mut body)?;
    let parser = Parser::new(&body);
    pulldown_cmark::html::push_html(&mut content, parser);
    post.content = content;
    return Ok(post);
}

#[cfg(test)]
mod tests {
    use super::*;
    use chrono::NaiveDate;

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
        let want: Vec<Tag> = vec!["en".to_string(), "Emacs".to_string(), "rant".to_string()];
        assert_eq!(got.tags, want);
    }
    #[test]
    fn test_read_tags_2() {
        let post = Post::new();
        let line = "## en";
        let got = read_tags(post, line).unwrap();
        let want: Vec<Tag> = vec!["en".to_string()];
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
        assert_eq!(post.tags, vec!["en".to_string(), "testing".to_string(),]);
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
