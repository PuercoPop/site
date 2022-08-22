package blog

import (
	"bytes"
	"context"
	"io/fs"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/civil"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
)

// TODO(javier): Initialize the blog from a embed.Fs
// - [ ] reverse-chronological index
// - [ ] an about page
// - [ ] taglists
// - [ ] Atom feed

type Site struct {
	fsys   fs.FS // TODO(javier): Should I use embed.FS instead?
	ByDate map[civil.Date]*Post
	ByTag  map[string]*Post
}

// Post represents a blog post written in markdown. Some metadata is embedded in
// the markdown as follows:
//
// - The first line should be a level 1 header which includes the title
// - The second line can be a level 3 header which would be the subtitle. This extension would ignore it.
// - The next line would be the tags of the post which is a comma delimited list inside a level 2 header.
// - The next line would be the date of the post which is in the format YYYY-MM-DD inside a level 2 header.
type Post struct {
	Title     string
	Tags      []string
	Published civil.Date
	Content   bytes.Buffer
}

const (
	STATE_TITLE = iota
	STATE_TAGS
	STATE_DATE
	STATE_DONE
)

// annotatePost walks the markdown document and copies any metadata found to the post.
func annotatePost(post *Post, data []byte) func(ast.Node, bool) (ast.WalkStatus, error) {
	state := STATE_TITLE
	return func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		switch state {
		case STATE_TITLE:
			if n.Kind() == ast.KindHeading && entering {
				hn := n.(*ast.Heading)
				if hn.Level == 1 {
					post.Title = string(n.Text(data))
					state = STATE_TAGS
				}
				return ast.WalkStatus(ast.WalkSkipChildren), nil
			}
		case STATE_TAGS:
			if n.Kind() == ast.KindHeading && entering {
				hn := n.(*ast.Heading)
				if hn.Level == 2 {
					post.Tags = strings.Split(string(n.Text(data)), ",")
				}
				return ast.WalkStatus(ast.WalkSkipChildren), nil
			}
		case STATE_DATE:
			if n.Kind() == ast.KindHeading && entering {
				hn := n.(*ast.Heading)
				if hn.Level == 2 {
					fmt := "2006-01-02"
					d, err := time.Parse(fmt, string(n.Text(data)))
					if err != nil {
						return ast.WalkStatus(ast.WalkStop), err
					}
					post.Published = civil.DateOf(d)
				}
				return ast.WalkStatus(ast.WalkSkipChildren), nil
			}
		case STATE_DONE:
			return ast.WalkStatus(ast.WalkStop), nil
		default:
			return ast.WalkStatus(ast.WalkContinue), nil
		}
		return ast.WalkStatus(ast.WalkContinue), nil

	}
}

// ReadPost reads a markdown file file and returns a Post.
func ReadPost(fpath string) (*Post, error) {
	data, err := os.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	post := &Post{}
	// How can we use goldmark parser here?
	reader := text.NewReader(data)
	md := goldmark.New(goldmark.WithRendererOptions(html.WithUnsafe()))
	doc := md.Parser().Parse(reader)
	visitor := annotatePost(post, data)
	// TODO(javier): Error check.
	ast.Walk(doc, visitor)
	// post.Content = md.Renderer().Render(d)

	// var buf bytes.Buffer
	// err = md.Convert(data, &buf)
	// if err != nil {
	// 	return nil, err
	// }
	// post.Content = buf
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
