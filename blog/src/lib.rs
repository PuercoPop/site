use axum::{
    extract::{Path as URLPath, State},
    response::Html,
    routing::get,
    Router,
};
use chrono::NaiveDate;
use minijinja::{context, Environment, Source};
use pulldown_cmark::{Event, Parser};
use regex::Regex;
use serde::Serialize;
use std::io::{self, BufRead, BufReader};
use std::path::Path;
use std::sync::Arc;
use std::{fs, io::Read};
use tokio_postgres::{Client, Error as PgError};

type Tag = String;

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

/// Holds to the global resources the application depends on
pub struct Context {
    /// The template engine
    templates: Environment<'static>,
    /// A connection to the database
    db: Client,
}

pub fn new_ctx(client: Client, template_dir: String) -> Context {
    let source = Source::from_path(template_dir);
    let mut env = Environment::new();
    env.set_source(source);

    Context {
        templates: env,
        db: client,
    }
}

/// Initializes the application. Takes the URL of the database to use.
pub fn new(ctx: Context) -> Router {
    // TODO: Can I remove the Arc wrapper?
    let ctx = Arc::new(ctx);
    let app = Router::new()
        .route("/", get(index))
        .route("/p/:slug", get(show_post))
        .with_state(ctx);
    app
}

#[axum::debug_handler]
async fn index(State(state): State<Arc<Context>>) -> Html<String> {
    let tmpl = state
        .templates
        .get_template("index.html")
        .expect("Unable to get template");
    let posts = recent_posts(&state.db).await.expect("IOU a ?");
    Html(
        tmpl.render(context!(latest_posts => posts))
            .expect("Unable to render template"),
    )
}

#[axum::debug_handler]
async fn show_post(
    State(state): State<Arc<Context>>,
    URLPath(slug): URLPath<String>,
) -> Html<String> {
    let tmpl = state
        .templates
        .get_template("post.html")
        .expect("Unable to get template");
    let post = post_by_slug(&state.db, slug).await.expect("IOU a ?");
    Html(
        tmpl.render(context!(post => post))
            .expect("Unable to render template"),
    )
}

static POST_BY_SLUG_QUERY: &str = "WITH posts AS (
  SELECT * FROM blog.posts WHERE slug = $1
), tags AS (
  SELECT post_id, array_agg(tag) AS tags FROM blog.post_tags
  WHERE post_id IN (SELECT post_id FROM posts) GROUP BY post_id
)
SELECT p.title, p.slug, p.draft, t.tags, p.published_at, p.content, p.path
 FROM posts p
NATURAL JOIN tags t
WHERE p.slug = $1";

async fn post_by_slug(client: &Client, slug: String) -> Result<Post, PgError> {
    let stmt = client.prepare(POST_BY_SLUG_QUERY).await?;
    let row = client.query_one(&stmt, &[&slug]).await?;
    let post = Post {
        slug: row.get("slug"),
        title: row.get("title"),
        pubdate: row.get("published_at"),
        tags: row.get("tags"),
        path: row.get("path"),
        draft: row.get("draft"),
        content: row.get("content"),
    };
    Ok(post)
}

static RECENT_POSTS_QUERY: &str = "WITH posts AS (
  select * from blog.posts order by published_at desc limit 5
), post_tags AS (
select post_id, array_agg(tag) as tags from blog.post_tags where post_id IN (select post_id from posts) group by post_id
)
select p.title, p.slug, p.draft, pt.tags, p.published_at, p.content, p.path from posts p natural join post_tags pt";

async fn recent_posts(client: &Client) -> Result<Vec<Post>, PgError> {
    let stmt = client.prepare(RECENT_POSTS_QUERY).await?;
    let posts: Vec<Post> = client
        .query(&stmt, &[])
        .await?
        .iter()
        .map(|row| Post {
            slug: row.get("slug"),
            title: row.get("title"),
            draft: row.get("draft"),
            tags: row.get("tags"),
            pubdate: row.get("published_at"),
            content: row.get("content"),
            path: row.get("path"),
        })
        .collect();
    Ok(posts)
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
