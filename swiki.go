package swiki

import (
	"net/http"

	"crawshaw.io/sqlite/sqlitex"
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

type store struct {
	pool *sqlitex.Pool
}

const DBPATH = "swkiki.db"

func NewStore(dbpath string) (*store, error) {
	pool, err := sqlitex.Open(dbpath, 0, 4)
	if err != nil {
		return nil, err
	}
	return &store{pool: pool}, nil

}
