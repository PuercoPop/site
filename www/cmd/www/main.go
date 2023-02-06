package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/PuercoPop/site/www"
)

var cert = flag.String("cert", "localhost+2.pem", "The certFile to use for http.ListenAndServeTLS")
var key = flag.String("key", "localhost+2-key.pem", "The keyFile to use for http.ListenAndServeTLS")

func main() {
	flag.Parse()
	handler := www.New(www.FSResources)
	log.Fatal(http.ListenAndServeTLS(":8080", *cert, *key, handler))
}
