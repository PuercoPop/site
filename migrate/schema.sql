-- -*- sql-product: postgres -*-
BEGIN;

CREATE SCHEMA blog IF NOT EXISTS;

CREATE SCHEMA finsta IF NOT EXISTS;

CREATE EXTENSION IF NOT EXISTS pg_stats;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE EXTENSION IF NOT EXISTS citext;

-- create public.migrations ()
-- blog

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

-- finsta
CREATE TABLE finsta.users (
  user_id integer PRIMARY KEY GENERATED always AS IDENTITY,
  email text NOT NULL UNIQUE, -- todo(javier): we want uniqueness to be case-insensitive
  password BYTEA NOT NULL,
  admin boolean NOT NULL DEFAULT FALSE
);

CREATE TABLE finsta.sessions (
  session_id bytea PRIMARY KEY,
  user_id integer NOT NULL,
  created_at time WITH time zone NOT NULL DEFAULT NOW(),
  expires_at time WITH time zone,
  FOREIGN KEY (user_id) REFERENCES users (user_id)
);

-- CREATE TABLE content_types (
--     content_type_id INTEGER PRIMARY KEY,
--     name TEXT NOT NULL UNIQUE
-- );
-- INSERT INTO content_types (name)
-- VALUES ('note'), ('photo'), ('article');
-- CREATE table posts (
--     page_id INTEGER PRIMARY KEY,
--     title TEXT NOT NULL, -- slug?
--     content TEXT NOT NULL,
--     published_at date not null default now(),
--     CREATED_AT DATETIME NOT NULL DEFAULT NOW()
--     -- author_id integer references
-- )

COMMIT;
