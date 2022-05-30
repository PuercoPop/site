package swiki

import (
	"context"
	"embed"
	"html/template"
	"log"
	"net/http"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type swiki struct {
	Mux *http.ServeMux
	T   *template.Template
}

//go:embed template/*.tmpl
var FSTemplates embed.FS

func New(dbpath string) *swiki {
	mux := http.NewServeMux()
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
	srv.Mux.HandleFunc("/sign-in/", srv.handleSignin())
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

func (srv *swiki) handleSignin(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			err := srv.T.ExecuteTemplate(w, "sign-in.html.tmpl", nil)
			// todo(javier): log error instead of dying.
			if err != nil {
				log.Fatalf("Could not render sign-in template. %s", err)
			}
		case http.MethodPost:
			// todo(javier): check credentials
			err := r.ParseForm()
			if err != nil {
				log.Fatalf("Could not parse form. %s", err)
			}
			email := req.PostForm.Get("email")
			password := req.PostForm.Get("password")
			var hash []byte
			err := db.QueryRow(context.TODO(), "SELECT password FROM users where email = $1", email).Scan(&hash)
			// If there is a user with that email.
			if errors.is(err, pgx.ErrNoRows) {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				// TODO(javier): Add an error to the template
				err := srv.T.ExecuteTemplate(w, "sign-in.html.tmpl", nil)
				// todo(javier): log error instead of dying.
				if err != nil {
					log.Fatalf("Could not render sign-in template. %s", err)
				}

			}
			err = bcrypt.CompareHashAndPassword(hash, password)
			if err != nil {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				// TODO(javier): Add an error to the template
				err := srv.T.ExecuteTemplate(w, "sign-in.html.tmpl", nil)
				// todo(javier): log error instead of dying.
				if err != nil {
					log.Fatalf("Could not render sign-in template. %s", err)
				}
			}
			// set cookie
			w.Write()

		}

	}
}

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
