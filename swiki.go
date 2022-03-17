package swiki

import (
	"embed"
	"html/template"
	"log"
	"net/http"

	"crawshaw.io/sqlite/sqlitex"
)

type swiki struct {
	Mux *http.ServeMux
	T   *template.Template
}

//go:embed template/*.tmpl
var FSTemplates embed.FS

func New(mux *http.ServeMux) *swiki {
	srv := &swiki{Mux: mux}
	srv.registerroutes()
	t, err := template.ParseFS(FSTemplates, "template/*.tmpl")
	if err != nil {
		log.Fatalf("Could not pare the templates: %s", err)
	}
	srv.T = t
	return srv
}

func (srv *swiki) registerroutes() {
	srv.Mux.HandleFunc("/", srv.indexFunc())
}

// Add an html/template here
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
