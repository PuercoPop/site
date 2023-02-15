-- -*- sql-product: postgres -*-
BEGIN;

CREATE SCHEMA blog IF NOT EXISTS;

CREATE TABLE blog.posts (
  post_id integer GENERATED always AS IDENTITY PRIMARY KEY,
  title text NOT NULL,
  slug text NOT NULL UNIQUE,
  path text NOT NULL UNIQUE,
  published_at time WITH time zone
);

CREATE INDEX IF NOT EXISTS blog.post_slugs ON blog.posts (slug);

CREATE INDEX IF NOT EXISTS blog.post_published_at ON blog.posts (published_at);

CREATE TABLE blog.tags (
  tag text NOT NULL PRIMARY KEY
);

CREATE TABLE blog.post_tags (
  post_id integer NOT NULL,
  tag text NOT NULL,
  FOREIGN KEY (post_id) REFERENCES blog.posts (post_id),
  FOREIGN KEY (tag) REFERENCES blog.tags (tag),
  PRIMARY KEY (post, tag)
);

COMMIT
