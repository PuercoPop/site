# Ergoproxy

We want:

- A process to hold any open connections while we update the services.
- To handle the renewal SSL certificates automatically.

Therefore we proxy.

- https://pkg.go.dev/golang.org/x/crypto/acme


## How it works

We start listening on a TCP socket. We start. We communicate the socket the
subprocess should though the environment variable `PROXY_SOCKET`. The Proxy
reads the first HTTP Line to determine the recipient URL, which is used to
decide which child process to send the request to. The proxy also sets the
`X-Forwarded-Host` header.

### Restarting a service

### Glossary/ Cast of characters / Dramatis Personae

- HTTP Proxy

- Sub-process, child process or services
