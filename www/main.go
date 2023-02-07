package main

import (
	"embed"
	"flag"
	"io/fs"
	"log"
	"net/http"
)

//go:embed resources/*.js resources/*.css resources/*.html
var FSResources embed.FS

type www struct {
	ResourceHandler http.Handler
}

func New(FSResources *embed.FS) *www {
	h := &www{}
	// TODO(javier): Move it to the main
	dir, err := fs.Sub(FSResources, "resources")
	if err != nil {
		log.Fatalf("Could not open resources directory: %s", err)
	}
	h.ResourceHandler = http.FileServer(http.FS(dir))
	return h
}

func (www *www) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	www.ResourceHandler.ServeHTTP(w, r)
}

var cert = flag.String("cert", "../ergoproxy/localhost+2.pem", "The certFile to use for http.ListenAndServeTLS")
var key = flag.String("key", "../ergoproxy/localhost+2-key.pem", "The keyFile to use for http.ListenAndServeTLS")

func main() {
	flag.Parse()
	handler := New(&FSResources)
	log.Fatal(http.ListenAndServeTLS(":8080", *cert, *key, handler))
}
