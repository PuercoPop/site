package swiki

import "net/http"

type swiki struct {
	mux *http.ServeMux
}

func New() *swiki {
	return &swiki{}
}

func NewMux() {
	mux := http.NewServeMux()
	mux.Handle("/", *indexFunc())
}

func indexFunc() *http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("Hello World"))
	}
}
