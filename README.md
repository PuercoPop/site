# site
## The code that powers my personal site.

It is inspired by the indieweb philosophy of self-hosting and ownership of
data. See [contributing.md](./docs/CONTRIBUTING.md) for information on how to
setup the project.

## Overview

The initial version of the site has the following components.

- ergo: A proxy whose main purpose is to upgrade other HTTP services w/o any
  downtime and keep the SSL certificates up to date..
- blog: Where I write my thoughts.
- finsta: An invite-only image board.
- webhookd: updates and redeploys other components upon learning of a new
  version.
- www: The landing page + webfinger endpoint.

For now all the components are written in Go, but that may change in the future.

# Development

## Setup

We can use docker for the database:

```shell
docker compose up -d postgres
```

Then we create the database for the entire site:

```shell
createdb -h localhost -U postgres site
```

Until we are live in production the we setup the database from a single
`schema.sql` file as opposed to migrations. Each project has a `sql/schema.sql`
file. TO setup the database schema we can use `psql -f`.. Ej. For the blog
project one can setup the database with:

```shell
psql -h locoalhost -U postgres -d site -f blog/sql/schema.sql
```
