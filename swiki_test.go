package swiki

import (
	"context"
	"testing"
)

func TestMain(m *testing.M) {
	// OpenDB collection
}

func teststore(t *testing.T) *Store {
	t.Helper()
	// TODO: Use :memory: database
	svc, err := NewStore("test_db")
	if err != nil {
		t.Fatalf("Could not open test db: %s", err)
	}
	return svc

}

func TestLatestPosts(t *testing.T) {
	store := teststore(t)
	t.Run("With no posts in the database", func(t *testing.T) {
		// setup
		// work
		// checks
		posts, err := store.ListRecentPosts(context.TODO(), 10)
		if err != nil {
			t.Errorf("Expected to return successfully. %s", err)
		}
		if len(posts) != 0 {
			t.Errorf("Expected an empty slice of Posts. Got %v", posts)
		}
	})
}

func TestReadPost(t *testing.T) {
	t.Run("Reads a post successfully", func(t *testing.T) {
		path := "./testdata/post.md"
		post, err := ReadPost(path)
		if err != nil {
			t.Errorf("Unexpected error returned. %s", err)
		}
		if post.Title != "Some title" {
			t.Errorf("Wrong title. got: %s. want: Some title", post.Title)
		}
	})
}
