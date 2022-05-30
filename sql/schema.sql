-- -*- sql-product: postgres -*-
BEGIN;

CREATE TABLE users (
       user_id INTEGER PRIMARY KEY generated always as identity,
       email TEXT NOT NULL UNIQUE, -- todo(javier): we want uniqueness to be case-insensitive
       password BYTEA NOT NULL,
       admin BOOLEAN NOT NULL DEFAULT false
);

create table sessions (
       session_id BYTEA PRIMARY KEY,
       user_id INTEGER NOT NULL,
       created_at TIME WITH TIME ZONE NOT NULL DEFAULT NOW(),
       expires_at TIME WITH TIME ZONE,
       foreign key (user_id) references users(user_id)
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
