package site

import (
	"embed"
	"net/http"

	"github.com/PuercoPop/site/blog"
)

// site is the top level handler
type site struct {
	WWW  *www
	Blog *blog.Site
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
		svc.WWW.ServeHTTP(w, r)
	}
}

// Add an html/template here

// func (srv *swiki) PageHandlerFunc() http.HandlerFunc {

// }
