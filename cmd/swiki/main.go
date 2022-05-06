package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/PuercoPop/swiki"
)

func main() {
	dbpath := flag.String("d", "swiki.db", "Path to th SQLite databse")
	flag.Parse()
	srv := swiki.New(dbpath)
	// TODO(javier): Obtain addr and cert files from flags
	log.Fatal(http.ListenAndServeTLS(":8080", "localhost+2.pem", "localhost+2-key.pem", srv.Mux))
}
