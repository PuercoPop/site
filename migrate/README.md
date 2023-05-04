# migrate
## A database migration tool

Runs any pending migrations and exists. The database migration scripts are
included in the executable, making it easy to deploy at the expense of flexibility.

```shell
$ ./migrate -d DBURL
```

# Principles

Give you the basic machinery but trust the developer knows what they are
doing. Keep things simple.

# Features

- Migrations only roll forward.
- The author of the migration is responsible for choosing and setting up the
  transactions they want to use.


# Build

```shell
go build ./...
```

# Usage

To apply any pending migrations we call `migrate` specifying the database.

```shell
migrate -d $dburl
```

By default it will the migrations that are embedded at build time from the
`migrations/` directory. Optionally we can specify a different directory to
obtain the list of migrations from. The migrations are always applied in
lexicographical order. We store which migrations have been applied in the
database. If they have already been applied they are skipped. We store both the
name of the migration and the checksum. If we encounter a migration that has
been applied with a different checksum we abort.

To dump the database schema we can use pg_dump, `pg_dump -s $dburl`.
