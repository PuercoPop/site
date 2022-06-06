package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/PuercoPop/site"
)

var dbpath = flag.String("d", "swiki.db", "Path to th SQLite databse")
var cert = flag.String("cert", "localhost+2.pem", "The certFile to use for http.ListenAndServeTLS")
var key = flag.String("key", "localhost+2-key.pem", "The keyFile to use for http.ListenAndServeTLS")

func main() {
	flag.Parse()
	srv := site.New(*dbpath)
	// TODO(javier): Obtain addr and cert files from flags
	log.Fatal(http.ListenAndServeTLS(":8080", *cert, *key, srv.Mux))
}
