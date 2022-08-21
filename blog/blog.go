package blog

import (
	"bytes"
	"cloud.google.com/go/civil"
	"context"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"os"
)

// TODO(javier): Initialize the blog from a embed.Fs
// - [ ] reverse-chronological index
// - [ ] an about page
// - [ ] taglists
// - [ ] Atom feed
import (
	"bytes"
	"context"
	"os"
	"io/fs"

	"cloud.google.com/go/civil"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

type Site struct{
	fsys fs.FS // TODO(javier): Should I use embed.FS instead?
	ByDate map[civil.Date]*Post
	ByTag map[string]*Post
}

type Post struct {
	title     string
	tags      []string
	published *civil.Date
	Content   bytes.Buffer
}

func ReadPost(fpath string) (*Post, error) {
	data, err := os.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	post := &Post{}
	// How can we use goldmark parser here?
	md := goldmark.New(goldmark.WithParserOptions(
		parser.WithHeadingAttribute(),
		parser.WithAutoHeadingID(),
		// parser.withASTTransformers(extractMetadata(&metadata)
	),
		goldmark.WithRendererOptions(html.WithUnsafe()))
	var buf bytes.Buffer
	err = md.Convert(data, &buf)
	if err != nil {
		return nil, err
	}
	post.Content = buf
	return post, nil
}

// Posts know how to render themselves as HTML
// func (p *Post)ServeHTTP(w httpResponseWriter, r *http.Request){}

type PostRepository interface {
	// Return the N most recent posts
	ListRecentPosts(ctx context.Context, n int) ([]*Post, error)
	Save(ctx context.Context, post Post) error
}

type PostMemRepository struct {
	posts []Post
}
