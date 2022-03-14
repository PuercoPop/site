package main

import (
	"log"

	"github.com/PuercoPop/swiki/swiki"
)

func main() {
	srv := swiki.New()
	// TODO: Replace for ListAndServeTLS,,,,nnm
	log.Fatal(srv.ListenAndServe())
}
