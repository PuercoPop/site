use::std::result::Result::Err;
use::std::path::Path;

#[derive(Debug)]
struct Post {
    title: String,
    pubdate: chrono::NaiveDate,
    draft: bool,
    tags: Vec<Tag>
}

#[derive(Debug)]
struct Tag {
    name: String
}

pub fn read_post(path: &Path) -> Result<Post, Err(String)> {
    return Ok(Post{})
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
        let post = read_post(path);
        assert_eq!(post.draft, true)
    }
    // fn test_draft_2() {
    //     path = Path::new("./testdata/post_01.md");
    //     post = read_post(path);
    //     assert_eq!(post.draft, false)
    // }
}
