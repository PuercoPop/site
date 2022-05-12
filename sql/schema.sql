-- -*- sql-product: sqlite -*-

PRAGMA foreign_keys = ON; -- connection specific setting?

create table users (
user_id integer primary key,
email text not null,
password text not null,
admin boolean not null default false
);

create table content_types(
    content_type_id integer primary key,
    name text not null unique
);
insert into content_types (name) values
('note'), ('photo'), ('article');

create table posts (
    page_id integer primary key,
    title text not null, -- slug?
    content text not null
    created_at datetime not null default NOW()
    -- author_id integer references
)
