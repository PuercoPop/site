package blog

import (
	"context"

	"cloud.google.com/go/civil"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RepositoryInterface interface {
	Upsert(ctx context.Context, post Post) error
	// Return the N most recent posts
	ListRecentPosts(ctx context.Context, n int) ([]*Post, error)
	FindBySlug(ctx context.Context, slug string) (*Post, error)
	ListByTag(ctx context.Context, tag string) ([]*Post, error)
	GroupByDate(ctx context.Context) ([]*ByDate, error)
}

type Repository struct {
	db *pgxpool.Pool
}

type ByDate struct {
	date  *civil.Date
	posts []*Post
}

var sqlRecentPosts = `
SELECT * FROM blog.posts LIMIT $1
`

// ListRecentPosts returns the N most recent posts
func (svc *Repository) ListRecentPosts(ctx context.Context, n int) ([]*Post, error) {
	rows, err := svc.db.Query(ctx, sqlRecentPosts, n)
	if err != nil {
		return nil, err
	}
	posts, err := pgx.CollectRows(rows, pgx.RowTo[*Post])
	if err != nil {
		return nil, err
	}
	return posts, nil
}
