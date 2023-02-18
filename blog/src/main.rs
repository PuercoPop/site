use pulldown_cmark::{Event, Parser};
use std::fs;
use std::io::{self, BufRead, BufReader};
use std::path::Path;

#[derive(Debug)]
pub struct Post {
    title: String,
    pubdate: chrono::NaiveDate,
    draft: bool,
    tags: Vec<Tag>, // path: String
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

fn read_title(line: &str) -> Result<String, io::Error> {
    let parser = Parser::new(line);
    for ev in parser {
        if let Event::Text(text) = ev {
            return Ok(text.to_string());
        }
    }
    Err(io::Error::from(io::ErrorKind::Other))
}

fn read_tags(line: &'static str) -> Result<Vec<Tag>, io::Error> {
    let parser = Parser::new(line);
    let mut ret: Vec<Tag> = Vec::new();
    for ev in parser {
        match ev {
            Event::Text(text) => {
                for tag in text.split(',') {
                    // trim
                    ret.push(Tag {
                        name: tag.trim().to_string(),
                    })
                }
                return Ok(ret);
            }
            _ => (),
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
        println!("line: {l:#?}")
    }

    // let input = fs::read_to_string(path).expect("Could not read file");
    // let parser = Parser::new(input.as_str());
    // for ev in parser {
    //     println!("ev: {:#?}", ev)
    // };
    let tags: Vec<Tag> = Vec::new();
    let pubdate = chrono::NaiveDate::from_ymd_opt(2023, 2, 15).unwrap();
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
        let line = "# Some title";
        let got = read_title(line).unwrap();
        assert_eq!(got, "Some title");
    }
    #[test]
    fn test_read_tags_1() {
        let line = "# en, Emacs, rant";
        let got = read_tags(line).unwrap();
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
        assert_eq!(got, want);
        // for (g, w) in got.iter().zip(want.iter_mut()) {
        //     assert_eq!(g, w)
        // }
    }
    #[test]
    fn test_read_tags_2() {
        let line = "# en";
        let got = read_tags(line).unwrap();
        let want: Vec<Tag> = vec![Tag {
            name: "en".to_string(),
        }];
        assert_eq!(got, want);
    }
    #[test]
    fn test_draft_1() {
        let path = Path::new("./testdata/draft_01.md");
        let post = read_post(path).expect("Could not read post");
        assert_eq!(post.draft, true)
    }
    #[test]
    fn test_draft_2() {
        let path = Path::new("./testdata/post_01.md");
        let post = read_post(path).expect("Could not read post");
        assert_eq!(post.draft, false)
    }
}
