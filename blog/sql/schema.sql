-- -*- sql-product: postgres -*-
BEGIN;

DROP SCHEMA IF EXISTS blog CASCADE;
CREATE SCHEMA IF NOT EXISTS blog ;

CREATE TABLE blog.posts (
  post_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  title TEXT NOT NULL,
  slug TEXT NOT NULL UNIQUE GENERATED ALWAYS AS (regexp_replace(title,' ', '-','g')) STORED,
  published_at DATE,
  draft BOOLEAN NOT NULL,
  content TEXT NOT NULL,
  path text NOT NULL UNIQUE
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
