package site

import (
	"embed"
	"net/http"
	"strings"

	"github.com/PuercoPop/site/blog"
)

// site is the top level handler
type site struct {
	// WWW *WWW
	// Blog *blog.Site
	WWW    http.Handler
	Blog   http.Handler
	Finsta http.Handler
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

func subdomain(r *http.Request) string {
	h := host(r)
	end := strings.Index(h, ".")
	return h[:end]
}

func (svc *site) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch h := subdomain(r); h {
	case "blog":
		svc.Blog.ServeHTTP(w, r)
	case "finsta":
		svc.Finsta.ServeHTTP(w, r)
	default:
		svc.WWW.ServeHTTP(w, r)
	}
}

// Add an html/template here

// func (srv *swiki) PageHandlerFunc() http.HandlerFunc {

// }
