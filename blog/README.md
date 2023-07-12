# Blog

## Local development

Update the post index

```shell
$ ./target/debug/import-blog -D ./content/posts/ -d postgres://postgres@localhost:5432/site
```

Serve the blog

```shell
$ ./target/debug/serve-blog -d postgres://postgres@localhost:5432/site -D ./templates
```
