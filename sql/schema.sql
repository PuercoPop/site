-- -*- sql-product: postgres -*-

CREATE TABLE users (
       user_id INTEGER PRIMARY KEY generated always as identity,
       email TEXT NOT NULL,
       password BYTEA NOT NULL,
       admin BOOLEAN NOT NULL DEFAULT false
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
