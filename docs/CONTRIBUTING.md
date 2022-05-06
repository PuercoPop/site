# Development
## Setting up local development.

swiki requires HTTPS. You can generate SSL ceritificages for localhost using [mkcert].

```shell
$ go install filippo.io/mkcert@latest
$ mkcert -install
$ mkcert localhost 127.0.0.1 ::1
```

[mkcert]: https://github.com/FiloSottile/mkcert

## Running the server locally

```shell
$ fd '.go$' . | entr -s 'go run ./cmd/swiki/main.go'
```
