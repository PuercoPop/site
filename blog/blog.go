package blog

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"path"
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
	BySlug map[string]*Post
	ByDate map[civil.Date][]*Post
	ByTag  map[string][]*Post
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
func ReadPost(fpath string, fsys fs.FS) (*Post, error) {
	data, err := fs.ReadFile(fsys, fpath)
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

func slugify(title string) string {
	return strings.ToLower(strings.ReplaceAll(title, " ", "-"))
}

// New initializes a new blog.
func New(blogFS fs.FS) *Site {
	site := &Site{}
	site.BySlug = make(map[string]*Post)
	site.ByTag = make(map[string][]*Post)
	site.ByDate = make(map[civil.Date][]*Post)
	// Walk the file-system for posts, loads them into memory and build an
	// index.
	fs.WalkDir(blogFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatalf("[blog.New]: %s", err)
		}
		// TODO(javier): Check it ends in markdown?
		if !d.IsDir() {
			post, err := ReadPost(path, blogFS)
			if err != nil {
				log.Fatalf("[blog.New]: %s", err)
			}
			for _, t := range post.Tags {
				xs := site.ByTag[t]
				site.ByTag[t] = append(xs, post)
			}
			ds := site.ByDate[post.Published]
			site.ByDate[post.Published] = append(ds, post)
			slug := slugify(post.Title)
			if site.BySlug[slug] != nil {
				log.Fatalf("Duplicated slug detected: %s", slug)
			}
			site.BySlug[slug] = post
		}
		return nil
	})
	return site
}

// ServeHTTP is the blogs entry point. The URL
// - /           -> The last five posts, with preview.
// - /tags       -> Renders an alphabetical list of tags.
// - /t/:tag     -> Renders an alphabetical list of posts tagged by :tag.
// - /p/:slug    -> Shows the post with :slug.
// - /archives/  -> A reverse chronological list of the posts.
// - /d/YYYY-M-D -> Renders an alphabetical list of posts published on YYYY-M-D.
// -> /atom.xml  -> The atom feed.
func (blog *Site) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)
	switch head {
	case "":
		blog.serveIndex(w, r)
	// case "tags":
	// 	serveTagList(w, r)
	// case "t":
	// 	serveTag(w, r)
	default:
		w.WriteHeader(404)

	}
}

// shiftPath splits the given path into the first segment (head) and
// the rest (tail). For example, "/foo/bar/baz" gives "foo", "/bar/baz".
// h/t: https://blog.merovius.de/posts/2017-06-18-how-not-to-use-an-http-router/
func shiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}

func (blog *Site) serveIndex(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	posts, err := blog.ListRecentPosts(ctx, 5)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.Write([]byte("<ol>"))
	for _, p := range posts {
		l := fmt.Sprintf("<li>%s published on %v</li>\n", p.Title, p.Published)
		w.Write([]byte(l))
	}
	w.Write([]byte("</ol>"))
}

// Posts know how to render themselves as HTML
// func (p *Post)ServeHTTP(w httpResponseWriter, r *http.Request){}

type PostRepository interface {
	// Return the N most recent posts
	ListRecentPosts(ctx context.Context, n int) ([]*Post, error)
	Save(ctx context.Context, post Post) error
}

func (blog *Site) ListRecentPosts(ctx context.Context, n int) ([]*Post, error) {
	var posts []*Post
	// TODO(javier): Replace with blog.ByDate.Keys once we have access to
	// generics.
	dates := make([]civil.Date, len(blog.ByDate))
	// We may need an array of dates with hits
	ix := 0
	for d := range blog.ByDate {
		// Does range start at 1?
		dates[ix] = d
	}
	// dates = sort.Sort(dates)
	postCount := 0
	for _, d := range dates {
		for _, post := range blog.ByDate[d] {
			posts = append(posts, post)
			postCount++
			if postCount >= n {
				return posts, nil
			}
		}
	}
	return posts, nil
}

type PostMemRepository struct {
	posts []Post
}
