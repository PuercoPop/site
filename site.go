package site

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/PuercoPop/site/blog"
	"github.com/jackc/pgx/v4/pgxpool"
)

// site is the top level handler
type site struct {
	Mux        *http.ServeMux
	Blog       *blog.Site
	t          *template.Template
	db         *pgxpool.Pool
	sessionsvc SessionService
}

//go:embed template/*.tmpl
var FSTemplates embed.FS

//go:embed resources/*.js resources/*.css
var FSResources embed.FS

//go:embed content/posts/*.md
var FSBlog embed.FS

func New() *site {
	h := &site{}
	h.Blog = blog.New(FSBlog)
	return h
}

func host(r *http.Request) string {
	if h := r.Header.Get("X-Forwarded-Host"); h != "" {
		return h
	}
	if h := r.Header.Get("Forwarded"); h != "" {
		return h
	}
	return r.Host
}

func (svc *site) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch h := host(r); h {
	case "blog":
		svc.Blog.ServeHTTP(w, r)
	default:
		// TODO(javier): Replace with "www" handler when we add one.
		svc.Blog.ServeHTTP(w, r)
	}
}

// func New(dbpath string) *site {
// 	h := &site{}
// 	t, err := template.ParseFS(FSTemplates, "template/*.tmpl")
// 	if err != nil {
// 		log.Fatalf("Could not pare the templates: %s", err)
// 	}
// 	h.t = t
// 	db, err := NewDB(context.Background(), dbpath)
// 	if err != nil {
// 		log.Fatalf("Could not connect to database: %s", err)
// 	}
// 	h.sessionsvc = &SessionStore{db: db}
// 	sm := &SessionMiddleware{svc: h.sessionsvc}
// 	h.Mux = http.NewServeMux()
// 	h.Mux.HandleFunc("/", sm.wrap(h.indexFunc()))
// 	h.Mux.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.FS(FSResources))))
// 	h.Mux.HandleFunc("/sign-in/", h.handleSignin())

// 	return h
// }

// Add an html/template here

// func (srv *swiki) PageHandlerFunc() http.HandlerFunc {

// }
