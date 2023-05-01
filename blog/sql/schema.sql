-- -*- sql-product: postgres -*-
BEGIN;

CREATE SCHEMA IF NOT EXISTS blog ;

CREATE TABLE blog.posts (
  post_id integer GENERATED always AS IDENTITY PRIMARY KEY,
  title text NOT NULL,
  slug text NOT NULL UNIQUE GENERATED ALWAYS AS (regexp_replace(title,' ', '-','g')) STORED,
  path text NOT NULL UNIQUE,
  published_at DATE
);

CREATE INDEX IF NOT EXISTS index_blog_posts_slugs ON blog.posts(slug);

CREATE INDEX IF NOT EXISTS index_blog_posts_published_at ON blog.posts(published_at);

CREATE TABLE blog.tags (
  tag text NOT NULL PRIMARY KEY
);

CREATE TABLE blog.post_tags (
  post_id integer NOT NULL,
  tag text NOT NULL,
  FOREIGN KEY (post_id) REFERENCES blog.posts (post_id),
  FOREIGN KEY (tag) REFERENCES blog.tags (tag),
  PRIMARY KEY (post_id, tag)
);

COMMIT;
