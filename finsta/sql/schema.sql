-- -*- sql-product: postgres -*-
begin;

DROP SCHEMA IF EXISTS finsta cascade;
CREATE SCHEMA IF NOT EXISTS finsta;

CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;

CREATE TABLE finsta.users (
  user_id INTEGER PRIMARY KEY GENERATED always AS IDENTITY,
  email TEXT NOT NULL UNIQUE, -- TODO: we want uniqueness to be case-insensitive
  password TEXT NOT NULL
  -- display_name text not null;
  -- name text not null;
);

CREATE TABLE finsta.admins (
       user_id INTEGER PRIMARY KEY,
       FOREIGN KEY (user_id) REFERENCES finsta.users (user_id)
);

CREATE TABLE finsta.sessions (
  session_id bytea PRIMARY KEY,
  user_id integer NOT NULL,
  created_at time WITH time zone NOT NULL DEFAULT NOW(),
  expires_at time WITH time zone,
  FOREIGN KEY (user_id) REFERENCES finsta.users (user_id)
);

CREATE TABLE finsta.known_media (
       user_id integer NOT NULL,
       checksum BYTEA NOT NULL,
       FOREIGN KEY (user_id) REFERENCES finsta.users (user_id),
       UNIQUE(user_id, checksum)
);

--- Use cases
-- 1. Upload an image to the archive of the user

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
create table finsta.media (
       checksum bytea not null
       -- TODO: Bucket URL? Is there some kind of structure we want to mimick?
);
COMMIT;
