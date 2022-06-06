package site

import (
	"embed"
	"html/template"
	"log"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
)

// site is the top level handler
type site struct {
	Mux *http.ServeMux
	t   *template.Template
	db  *pgxpool.Pool
}

//go:embed template/*.tmpl
var FSTemplates embed.FS

func New(dbpath string) *site {
	h := &site{}
	t, err := template.ParseFS(FSTemplates, "template/*.tmpl")
	if err != nil {
		log.Fatalf("Could not pare the templates: %s", err)
	}
	h.t = t
	h.Mux = http.NewServeMux()
	h.Mux.HandleFunc("/", h.indexFunc())
	h.Mux.HandleFunc("/sign-in/", h.handleSignin())

	return h
}

// Add an html/template here

// func (srv *swiki) PageHandlerFunc() http.HandlerFunc {

// }

// type Store struct {
// 	pool *sqlitex.Pool
// }

// const DBPATH = "swkiki.db"

// func NewStore(dbpath string) (*Store, error) {
// 	pool, err := sqlitex.Open(dbpath, 0, 4)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &Store{pool: pool}, nil

// }

type Post struct {
	Title   string
	Content string
	// published *civil.Date
}

// func ReadPost(path string) (*Post, error) {
// 	return &Post{}, nil
// }

// // Posts know how to render themselves as HTML
// // func (p *Post)ServeHTTP(w httpResponseWriter, r *http.Request){}

// type PostsDBAL interface {
// 	// Return the N most recent posts
// 	ListRecentPosts(ctx context.Context, n int) ([]*Post, error)
// 	Save(ctx context.Context, post Post) error
// }

// const sqlListRecentPosts = `
// SELECT * FROM POSTS LIMIT $1`

// func (svc *Store) ListRecentPosts(ctx context.Context, n int64) ([]*Post, error) {
// 	conn := svc.pool.Get(ctx)
// 	defer svc.pool.Put(conn)
// 	stmt, err := conn.Prepare(sqlListRecentPosts)
// 	if err != nil {
// 		return nil, fmt.Errorf("Could not prepare the query %w", err)
// 	}
// 	stmt.BindInt64(1, n)
// 	return nil, fmt.Errorf("iou")
// }
