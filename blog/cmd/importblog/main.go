package main

import (
	"context"
	"flag"
	"io/fs"
	"log"

	"github.com/PuercoPop/site"
	"github.com/PuercoPop/site/blog"
	"github.com/jackc/pgx/v5/pgxpool"
)

var dburl = flag.String("d", "swiki.db", "The URL where the database is located.")
var blogdir = flag.String("D", "", "The directory where the blog posts are.")

func main() {
	ctx := context.Background()
	flag.Parse()
	pool, err := pgxpool.New(ctx, *dburl)
	if err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}
	svc := blog.NewPGRespository(pool)
	FSBlog := site.FSBlog
	fs.WalkDir(FSBlog, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatalf("[main.WalkDir]: %s", err)
		}
		if !d.IsDir() {
			post, err := blog.ReadPost(path, FSBlog)
			if err != nil {
				log.Fatalf("Could not read post at %s: %s", path, err)
			}
			err = svc.Import(ctx, post)
			if err != nil {
				log.Fatalf("Could not import post: %s", err)
			}

		}
		return nil

	})
}
