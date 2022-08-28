package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/PuercoPop/site"
)

var dbpath = flag.String("d", "swiki.db", "Path to th SQLite databse")
var cert = flag.String("cert", "localhost+2.pem", "The certFile to use for http.ListenAndServeTLS")
var key = flag.String("key", "localhost+2-key.pem", "The keyFile to use for http.ListenAndServeTLS")

func main() {
	flag.Parse()
	go func() {
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
		<-sc
		os.Exit(0)
	}()
	srv := site.New()
	// TODO(javier): Obtain addr and cert files from flags
	log.Fatal(http.ListenAndServeTLS(":8080", *cert, *key, srv))
}
