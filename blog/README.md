# Blog

## Post Format

Posts are written in markdown. Instead of adding a front-matter the information
is embedded in the markdown. The first line of the post should be a level 1
heard that is the title. The next line can be a level 2 header, in which case it
is the sub-title. Following that we would encounter a level 3 header which
contains the tags as a comma separated list and then the date in YYYY-MM-DD
format.

## Local development

Update the post index

```shell
$ ./target/debug/import-blog -D ./content/posts/ -d postgres://postgres@localhost:5432/site
```

Serve the blog

```shell
$ ./target/debug/serve-blog -d postgres://postgres@localhost:5432/site -D ./templates
```
