package site

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Post struct {
	Title   string
	Content string
	// published *civil.Date
}

// func ReadPost(path string) (*Post, error) {
// 	return &Post{}, nil
// }

// Posts know how to render themselves as HTML
// func (p *Post)ServeHTTP(w httpResponseWriter, r *http.Request){}

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

func (svc *PostStore) Save(ctx context.Context, post Post) error {
	return nil
}
