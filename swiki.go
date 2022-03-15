package swiki

import "net/http"

type swiki struct {
	*http.ServeMux
}

func New() *swiki {
	srv := &swiki{}
	srv.registerroutes()
	return srv
}

// func NewMux() {
// 	mux := http.NewServeMux()
// 	mux.Handle("/", indexFunc())
// }

func (srv *swiki) registerroutes() {
	srv.Handle("/", srv.indexFunc())
}

func (srv *swiki) indexFunc() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("Hello World"))
	}
}
