# Ergoproxy

We want to upgrade services without any downtime. To do so we need a process to
hold the connections open. Therefore we proxy.

## Usage

```shell
$ ergoproxy --subdomain www --service ./www-srv --host localhost --port 80
```

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

## References

- https://copyconstruct.medium.com/file-descriptor-transfer-over-unix-domain-sockets-dcbbf5b3b6ec
- Zero Downtime Release:Disruption-free Load Balancing of a Multi-Billion User
  Website
  https://research.facebook.com/publications/zero-downtime-release-disruption-free-load-balancing-of-a-multi-billion-user-website/
  https://dl.acm.org/doi/abs/10.1145/3387514.3405885
- https://github.com/sozu-proxy/sozu
