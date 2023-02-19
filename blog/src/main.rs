use chrono::NaiveDate;
use pulldown_cmark::{Event, Parser};
use regex::Regex;
use std::fs;
use std::io::{self, BufRead, BufReader};
use std::path::Path;

#[derive(Debug, Default)]
pub struct Post {
    title: String,
    pubdate: NaiveDate,
    draft: bool,
    tags: Vec<Tag>,
    // path: String
}
impl Post {
    pub(crate) fn new() -> Post {
        Post::default()
    }
}

#[derive(Debug, PartialEq)]
struct Tag {
    name: String,
}

enum MetadataParseState {
    TitleLine,
    Tags,
    DateLine,
    End, // End of front matter
}

type FSM = MetadataParseState;

// Error return when read_post fails, explaining why it wasn't possible to read the post.
// type PostParseError {
//     std::io::Error,
//     BadFormat // Or maybe use std::io::Error::other instead
// }

fn read_title(mut post: Post, line: &str) -> Result<Post, io::Error> {
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
    Err(io::Error::from(io::ErrorKind::Other))
}

fn read_tags(mut post: Post, line: &'static str) -> Result<Post, io::Error> {
    let parser = Parser::new(line);
    let mut tags: Vec<Tag> = Vec::new();
    for ev in parser {
        match ev {
            Event::Text(text) => {
                for tag in text.split(',') {
                    // trim
                    tags.push(Tag {
                        name: tag.trim().to_string(),
                    })
                }
                post.tags = tags;
                return Ok(post);
            }
            _ => (),
        }
    }
    Err(io::Error::from(io::ErrorKind::Other))
}

fn read_pubdate(mut post: Post, line: &'static str) -> Result<Post, io::Error> {
    let parser = Parser::new(line);
    for ev in parser {
        if let Event::Text(text) = ev {
            let pubdate = NaiveDate::parse_from_str(&text, "%Y-%m-%d").expect("I should use ?");
            post.pubdate = pubdate;
            return Ok(post);
        }
    }
    Err(io::Error::from(io::ErrorKind::Other))
}

// Reads the meta-data embedded in the markdown document and returns a Post.
pub fn read_post(path: &Path) -> Result<Post, ()> {
    let fd = fs::File::open(path).expect("Could not open file");
    let reader = BufReader::new(fd);
    let _state = FSM::TitleLine;
    for line in reader.lines() {
        let l = line.expect("Could not extract line contents");
        // match state {
        //     FSM::TitleLine => read_title(l),
        //     FSM::Tags => read_tag(l),
        //     FSM::DateLine => read_date(l),
        //     FSM::End => break;
        // }
        // We need to add a case using the state enum

        // TODO(javier): Handle EOF?
        if l.is_empty() {
            println!("End of front-matter")
        }
        // println!("line: {l:#?}")
    }

    // let input = fs::read_to_string(path).expect("Could not read file");
    // let parser = Parser::new(input.as_str());
    // for ev in parser {
    //     println!("ev: {:#?}", ev)
    // };
    let tags: Vec<Tag> = Vec::new();
    let pubdate = NaiveDate::from_ymd_opt(2023, 2, 15).unwrap();
    Ok(Post {
        title: "".to_string(),
        pubdate,
        draft: true,
        tags,
    })
}

fn main() {
    println!("Hello, world!");
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
    #[ignore]
    fn test_read_pubdate_3() {
        let post = Post::new();
        let line = "## 2022-2-31"; // impossible date
        let got = read_pubdate(post, line).unwrap();
        let want = NaiveDate::from_ymd_opt(2022, 2, 15).unwrap();
        assert_eq!(got.pubdate, want);
    }
    // TODO(javier): Move this tests to integration
    #[test]
    fn test_integration_1() {
        let path = Path::new("./testdata/draft_01.md");
        let post = read_post(path).expect("Could not read post");
        assert_eq!(post.draft, true)
    }
    #[test]
    fn test_integration_2() {
        let path = Path::new("./testdata/post_01.md");
        let post = read_post(path).expect("Could not read post");
        assert_eq!(post.draft, false)
    }
}
