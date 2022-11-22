package www

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/PuercoPop/site/blog"
	"github.com/PuercoPop/site/finsta"
)

type www struct {
	t               *template.Template
	ResourceHandler http.Handler
}

func New(dbpath string, FSResources *embed.FS, FSTemplates embed.FS) (*www, error) {
	h := &www{}
	h.ResourceHandler = http.FileServer(http.FS(FSResources))
	t, err := template.ParseFS(FSTemplates, "template/*.tmpl")
	if err != nil {
		return nil, fmt.Errorf("Could not pare the templates: %w", err)
	}
	h.t = t
	if err != nil {
		return nil, fmt.Errorf("Could not connect to database: %w", err)
	}
	return h, nil
}

func (www *www) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)
	switch head {
	case "resources":
		www.ResourceHandler.ServeHTTP(w, r)
	default:
		www.serveIndex(w, r)
	}
}

func (h *www) serveIndex(w http.ResponseWriter, r *http.Request) {
	type data struct {
		LatestPosts []blog.Post
		CurrentUser *finsta.User
	}
	posts := []blog.Post{
		{
			Title:   "Awesome Post!",
			Content: bytes.NewBufferString("lololol"),
		},
	}
	d := data{LatestPosts: posts, CurrentUser: nil}
	err := h.t.ExecuteTemplate(w, "index.html.tmpl", d)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error rendering tempalte. %s", err)
		return
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
