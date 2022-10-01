package www

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/PuercoPop/site"
	"github.com/PuercoPop/site/blog"
	"github.com/PuercoPop/site/database"
	"github.com/jackc/pgx/v4/pgxpool"
)

type www struct {
	t               *template.Template
	db              *pgxpool.Pool
	sm              *site.SessionMiddleware
	sessionsvc      site.SessionService
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
	db, err := database.New(context.Background(), dbpath)
	if err != nil {
		return nil, fmt.Errorf("Could not connect to database: %w", err)
	}
	h.sessionsvc = &site.SessionStore{db: db}
	h.sm = &site.SessionMiddleware{svc: h.sessionsvc}
	return h, nil
}

func (www *www) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)
	switch head {
	case "":
		www.serveIndex(w, r)
	case "resources":
		www.ResourcesHandler.ServeHTTP(w, r)
	case "sign-in":
		www.serveSignIn(w, r)

	}
}

func (h *www) serveIndex(w http.ResponseWriter, r *http.Request) {
	type data struct {
		LatestPosts []blog.Post
		CurrentUser *site.User
	}
	if err := h.sm.agument; err != nil {
		log.Fatalf(err)
		return
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

func (h *www) serveSignIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err := h.t.ExecuteTemplate(w, "sign-in.html.tmpl", nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Error rendering tempalte. %s", err)
			return
		}
	case http.MethodPost:
		// todo(javier): check credentials
		err := r.ParseForm()
		if err != nil {
			log.Fatalf("Could not parse form. %s", err)
		}
		email := r.PostForm.Get("email")
		password := r.PostForm.Get("password")
		sid, err := h.sessionsvc.New(ctx, email, password)
		if err != nil {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			// TODO(javier): Add an error to the template
			err := h.t.ExecuteTemplate(w, "sign-in.html.tmpl", nil)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Printf("Could not render sign-in template. %s\n", err)
			}
		}
		cookie := &http.Cookie{Name: "sid", Value: string(sid)}
		http.SetCookie(w, cookie)
		url := r.Form.Get("redirect_to")
		if url == "" {
			url = "/"
		}
		http.Redirect(w, r, url, http.StatusSeeOther)
		w.Write([]byte("login successful."))
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
