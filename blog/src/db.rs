use crate::post::{Post, Tag};
use serde::Serialize;
use tokio_postgres::{Client, Error as PgError};

static INSERT_POST_QUERY: &str = "INSERT INTO blog.posts (title, published_at, draft, content, path)
VALUES ($1, $2, $3, $4, $5) ON CONFLICT (path) DO UPDATE SET
title = EXCLUDED.title, published_at = EXCLUDED.published_at, draft = EXCLUDED.draft, content = EXCLUDED.content
RETURNING post_id";
static INSERT_TAGS_QUERY: &str = "INSERT INTO blog.tags (TAG) VALUES ($1) ON CONFLICT DO NOTHING";
static INSERT_POST_TAGS_QUERY: &str =
    "INSERT INTO blog.post_tags (post_id, tag) VALUES ($1, $2) ON CONFLICT DO NOTHING";

pub async fn store_post(client: &Client, post: Post) -> Result<(), PgError> {
    let stmt = client.prepare(INSERT_POST_QUERY).await?;
    let rows = client
        .query(
            &stmt,
            &[
                &post.title,
                &post.pubdate,
                &post.draft,
                &post.content,
                &post.path,
            ],
        )
        .await?;
    let post_id: i32 = rows.iter().next().unwrap().get(0);
    let stmt = client
        .prepare(INSERT_TAGS_QUERY)
        .await
        .expect("Could not prepare second statement");
    let insert_post_tags_stmt = client.prepare(INSERT_POST_TAGS_QUERY).await.unwrap();
    for tag in post.tags {
        let _ret = client
            .query(&stmt, &[&tag])
            .await
            .expect("Could not insert tag");
        let _ret = client
            .query(&insert_post_tags_stmt, &[&post_id, &tag])
            .await
            .unwrap();
    }
    Ok(())
}

static REMOVE_UNUSED_POST_TAGS_QUERY: &str = "WITH unused AS (
SELECT tag FROM blog.tags NATURAL LEFT OUTER JOIN blog.post_tags WHERE post_id IS NULL
)
DELETE FROM blog.tags WHERE tag IN (SELECT tag from unused)";

pub async fn remove_unused_tags(client: &Client) -> Result<(), PgError> {
    let stmt = client.prepare(REMOVE_UNUSED_POST_TAGS_QUERY).await?;
    let _ret = client.query(&stmt, &[]).await?;
    Ok(())
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

pub(crate) async fn post_by_slug(client: &Client, slug: String) -> Result<Post, PgError> {
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
select p.title, p.slug, p.draft, pt.tags, p.published_at, p.content, p.path from posts p natural join post_tags pt
order by p.published_at desc";

pub(crate) async fn recent_posts(client: &Client) -> Result<Vec<Post>, PgError> {
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

static TAGS_COUNT_QUERY: &str =
    "SELECT tag, count(*) as count FROM blog.post_tags GROUP BY tag ORDER BY tag";

#[derive(Serialize)]
pub(crate) struct TagEntry {
    name: Tag,
    count: i64,
}

pub(crate) async fn tags_count(client: &Client) -> Result<Vec<TagEntry>, PgError> {
    let stmt = client.prepare(TAGS_COUNT_QUERY).await?;
    let tags: Vec<TagEntry> = client
        .query(&stmt, &[])
        .await?
        .iter()
        .map(|row| TagEntry {
            name: row.get("tag"),
            count: row.get("count"),
        })
        .collect();
    Ok(tags)
}

static POSTS_BY_TAG_QUERY: &str = "WITH posts AS (
  SELECT p.* FROM blog.post_tags pt
  NATURAL JOIN blog.posts p
  WHERE pt.tag = $1
), tags AS (
  SELECT post_id, array_agg(tag) AS tags
  FROM blog.post_tags
  WHERE post_id IN (SELECT post_id FROM posts)
  GROUP BY post_id
)
SELECT posts.*, t.tags
FROM posts NATURAL JOIN tags t
ORDER BY posts.published_at DESC";

/// Returns all the posts tagged by `tag'
pub(crate) async fn posts_by_tag(client: &Client, tag: &Tag) -> Result<Vec<Post>, PgError> {
    let stmt = client.prepare(POSTS_BY_TAG_QUERY).await?;
    let posts: Vec<Post> = client
        .query(&stmt, &[tag])
        .await?
        .iter()
        .map(|row| Post {
            slug: row.get("slug"),
            title: row.get("title"),
            draft: row.get("draft"),
            tags: row.get("tags"),
            content: "".to_string(),
            pubdate: row.get("published_at"),
            path: row.get("path"),
        })
        .collect();
    Ok(posts)
}

static ALL_POSTS_QUERY: &str = "WITH posts AS (
  select * from blog.posts order by published_at desc
), post_tags AS (
select post_id, array_agg(tag) as tags from blog.post_tags where post_id IN (select post_id from posts) group by post_id
)
select p.title, p.slug, p.draft, pt.tags, p.published_at, p.content, p.path from posts p natural join post_tags pt
order by p.published_at desc";

// TODO: Once we move to a paginated feed we can remove this method and used recent_posts directly instead.
pub(crate) async fn all_posts(client: &Client) -> Result<Vec<Post>, PgError> {
    let stmt = client.prepare(ALL_POSTS_QUERY).await?;
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
