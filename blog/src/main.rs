// use::std::result::Result::Err;
use::std::path::Path;

#[derive(Debug)]
pub struct Post {
    title: String,
    pubdate: chrono::NaiveDate,
    draft: bool,
    tags: Vec<Tag>
}

#[derive(Debug)]
struct Tag {
    name: String
}

pub fn read_post(path: &Path) -> Result<Post, ()> {
    let tags: Vec<Tag> = Vec::new();
    let pubdate = chrono::NaiveDate::from_ymd_opt(2023, 2, 15).unwrap();
    return Ok(Post{title: "".to_string(), pubdate: pubdate, draft: true, tags: tags})
}

fn main() {
    println!("Hello, world!");
}

#[cfg(test)]
mod tests {
    use super::*;
    #[test]
    fn test_draft_1() {
        let path = Path::new("./testdata/draft_01.md");
        let post = read_post(path).expect("Could not read post");
        assert_eq!(post.draft, true)
    }
    #[test]
    fn test_draft_2() {
        let path = Path::new("./testdata/post_01.md");
        let post = read_post(path).expect("Could not read post");;
        assert_eq!(post.draft, false)
    }
}
