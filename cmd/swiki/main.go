package main

import (
	"log"
	"net/http"

	"github.com/PuercoPop/swiki"
)

func main() {
	mux := http.NewServeMux()
	srv := swiki.New(mux)
	// TODO(javier): Obtain addr and cert files from flags
	log.Fatal(http.ListenAndServeTLS(":8080", "localhost+2.pem", "localhost+2-key.pem", srv.Mux))
}
