package blog

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"log"
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
	ByTitle map[string]*Post // TODO(javier: probably should be slug instead)
	ByDate  map[civil.Date][]Post
	ByTag   map[string][]Post
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
	Content   *bytes.Buffer
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
				return ast.WalkSkipChildren, nil
			}
		case STATE_TAGS:
			if n.Kind() == ast.KindHeading && entering {
				hn := n.(*ast.Heading)
				if hn.Level == 2 {
					tags := strings.Split(string(n.Text(data)), ",")
					for ix := range tags {
						tags[ix] = strings.TrimSpace(tags[ix])
					}
					post.Tags = tags
					state = STATE_DATE
				}
				return ast.WalkSkipChildren, nil
			}
		case STATE_DATE:
			if n.Kind() == ast.KindHeading && entering {
				hn := n.(*ast.Heading)
				if hn.Level == 2 {
					timefmt := "2006-1-2"
					d, err := time.Parse(timefmt, string(n.Text(data)))
					if err != nil {
						return ast.WalkStop, err
					}
					post.Published = civil.DateOf(d)
					state = STATE_DONE
				}
				return ast.WalkSkipChildren, nil
			}
		case STATE_DONE:
			return ast.WalkStop, nil
		default:
			return ast.WalkContinue, nil
		}
		return ast.WalkContinue, nil

	}
}

// ReadPost reads a markdown file file and returns a Post.
func ReadPost(fpath string) (*Post, error) {
	data, err := os.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	post := &Post{}
	reader := text.NewReader(data)
	md := goldmark.New(goldmark.WithRendererOptions(html.WithUnsafe()))
	doc := md.Parser().Parse(reader)
	visitor := annotatePost(post, data)
	err = ast.Walk(doc, visitor)
	if err != nil {
		return nil, err
	}
	// TODO(javier): Save rendered markdown in Content field.
	var buf bytes.Buffer
	err = md.Renderer().Render(&buf, data, doc)
	if err != nil {
		return nil, err
	}
	post.Content = &buf
	return post, nil
}

// New initializes a new blog.
func New(blogFS fs.FS) *Site {
	site := &Site{}
	site.ByTag = make(map[string][]Post)
	// TODO(javier): Walk the file-system for posts, loads them into memory
	// and build an index.
	// TODO(javier): Replace the testdata with .
	fs.WalkDir(blogFS, "testdata", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatalf("[blog.New]: %s", err)
		}
		// TODO Check it ends in markdown?
		if !d.IsDir() {
			post, err := ReadPost(path)
			if err != nil {
				log.Fatalf("[blog.New]: %s", err)
			}
			for ix, t := range post.Tags {
				xs := site.ByTag[t]
				site.ByTag[t] = append(xs, post)
			}
			fmt.Printf("=> %s\n", post.Title)

		}
		return nil
	})
	return site
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
