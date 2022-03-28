package swiki

import (
	"context"
	"embed"
	"fmt"
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
		posts := []Post{
			{
				Title:   "Awesome Post!",
				Content: "lololol",
			},
		}
		data := struct{ LatestPosts []Post }{LatestPosts: posts}
		srv.T.ExecuteTemplate(res, "index.html.tmpl", data)
		// How to check an error here
	}
}

type Store struct {
	pool *sqlitex.Pool
}

const DBPATH = "swkiki.db"

func NewStore(dbpath string) (*Store, error) {
	pool, err := sqlitex.Open(dbpath, 0, 4)
	if err != nil {
		return nil, err
	}
	return &Store{pool: pool}, nil

}

type Post struct {
	Title   string
	Content string
	// published *civil.Date
}

type PostsDBAL interface {
	// Return the N most recent posts
	ListRecentPosts(ctx context.Context, n int) ([]*Post, error)
	Save(ctx context.Context, post Post) error
}

const sqlListRecentPosts = `
SELECT * FROM POSTS LIMIT $1`

func (svc *Store) ListRecentPosts(ctx context.Context, n int64) ([]*Post, error) {
	conn := svc.pool.Get(ctx)
	defer svc.pool.Put(conn)
	stmt, err := conn.Prepare(sqlListRecentPosts)
	if err != nil {
		return nil, fmt.Errorf("Could not prepare the query %w", err)
	}
	stmt.BindInt64(1, n)
	return nil, fmt.Errorf("iou")
}
