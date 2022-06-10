package site

import (
	"context"
	"embed"
	"fmt"
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
	h.Mux = http.NewServeMux()
	h.Mux.HandleFunc("/", h.indexFunc())
	h.Mux.HandleFunc("/sign-in/", h.handleSignin())

	return h
}

// Add an html/template here

// func (srv *swiki) PageHandlerFunc() http.HandlerFunc {

// }

type User struct {
	Email string
}

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

type PostService interface {
	// Return the N most recent posts
	ListRecentPosts(ctx context.Context, n int) ([]*Post, error)
	Save(ctx context.Context, post Post) error
}

type PostStore struct {
	db *pgxpool.Pool
}

const sqlListRecentPosts = `SELECT * FROM POSTS LIMIT $1`

func (svc *PostStore) ListRecentPosts(ctx context.Context, n int64) ([]*Post, error) {
	rows, err := svc.db.Query(ctx, sqlListRecentPosts, n)
	if err != nil {
		return nil, fmt.Errorf("Could not prepare the query %w", err)
	}
	defer rows.Close()
	for rows.Next() {

	}
	return nil, fmt.Errorf("iou")
}
