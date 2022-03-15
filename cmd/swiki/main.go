package main

import (
	"log"

	"github.com/PuercoPop/swiki"
)

func main() {
	srv := swiki.New()
	// TODO(javier): Obtain addr and cert files from flags
	log.Fatal(srv.ListenAndServeTLS(":8080", "localhost+2.pem", "localhost+2-key.pem", srv))
}
