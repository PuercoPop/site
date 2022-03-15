package swiki

import (
	"net/http"
)

type swiki struct {
	Mux *http.ServeMux
}

func New(mux *http.ServeMux) *swiki {
	srv := &swiki{Mux: mux}
	srv.registerroutes()
	return srv
}

func (srv *swiki) registerroutes() {
	srv.Mux.HandleFunc("/", srv.indexFunc())
}

func (srv *swiki) indexFunc() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("Hello World"))
	}
}
