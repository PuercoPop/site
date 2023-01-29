# migrate
## A database migration tool

Runs any pending migrations and exists. The database migration scripts are
included in the executable, making it easy to deploy at the expense of flexibility.

```shell
$ ./migrate -d DBURL
```

# Build

```shell
go build ./...
```
