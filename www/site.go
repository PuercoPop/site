package www

import (
	"embed"
	"net/http"
)

//go:embed resources/*.js resources/*.css resources/*.html
var FSResources embed.FS

type www struct {
	ResourceHandler http.Handler
}

func New(FSResources *embed.FS) *www {
	h := &www{}
	h.ResourceHandler = http.FileServer(http.FS(FSResources))
	return h
}

func (www *www) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	www.ResourceHandler.ServeHTTP(w, r)
}
