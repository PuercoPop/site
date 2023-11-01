# Finsta

## Overview

- Each user has a collection of private media, an archive, and a gallery where
  they can share the media. The gallery can be organized into albums (rooms?,
  expositions?).

## Development setup

```shell
createdb finsta -h localhost -U postgres
psql -h localhost -U postgres -d finsta -f sql/schema.sql
psql -h localhost -U postgres -d finsta -f sql/seeds.sql
```

## Affordances

As a User I need to be able to:

- Login to the site.
- Share pictures
- Effects?
- Upload by default in private/backup media from the phone.

## Companion App

[finsta-droid](https://github.com/PuercoPop/finsta-droid)
