package blog

import (
	"context"

	"cloud.google.com/go/civil"
)


type Post struct {
	title     string
	tags      []string
	published *civil.Date
	content   string
}

func ReadPost(path string) (*Post, error) {
	return &Post{}, nil
}

// Posts know how to render themselves as HTML
// func (p *Post)ServeHTTP(w httpResponseWriter, r *http.Request){}

type PostRepository interface {
	// Return the N most recent posts
	ListRecentPosts(ctx context.Context, n int) ([]*Post, error)
	Save(ctx context.Context, post Post) error
}

type PostMemRepository struct {
	posts []Post
}
