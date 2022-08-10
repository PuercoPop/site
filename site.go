package site

import (
	"context"
	"embed"
	"html/template"
	"log"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
)

// site is the top level handler
type site struct {
	Mux        *http.ServeMux
	t          *template.Template
	db         *pgxpool.Pool
	sessionsvc SessionService
}

//go:embed template/*.tmpl
var FSTemplates embed.FS

//go:embed resources/*.[js|css]
var FSResources embed.FS

func New(dbpath string) *site {
	h := &site{}
	t, err := template.ParseFS(FSTemplates, "template/*.tmpl")
	if err != nil {
		log.Fatalf("Could not pare the templates: %s", err)
	}
	h.t = t
	db, err := NewDB(context.Background(), dbpath)
	if err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}
	h.sessionsvc = &SessionStore{db: db}
	sm := &SessionMiddleware{svc: h.sessionsvc}
	h.Mux = http.NewServeMux()
	h.Mux.HandleFunc("/", sm.wrap(h.indexFunc()))
	h.Mux.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.FS(FSResources))))
	h.Mux.HandleFunc("/sign-in/", h.handleSignin())

	return h
}

// Add an html/template here

// func (srv *swiki) PageHandlerFunc() http.HandlerFunc {

// }
