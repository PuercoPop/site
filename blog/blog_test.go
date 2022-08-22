package blog

import (
	"bytes"
	"testing"
	"time"

	"cloud.google.com/go/civil"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestReadPost(t *testing.T) {
	tt := []struct {
		name string
		path string
		want *Post
	}{{
		name: "read post 1",
		path: "testdata/post.md",
		want: &Post{Title: "Some title",
			Tags:      []string{"en", "testing"},
			Published: civil.Date{Year: 2022, Month: time.March, Day: 30},
			Content:   bytes.NewBufferString(""),
		},
	}}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ReadPost(tc.path)
			if err != nil {
				t.Errorf("Could not read post successfully. %s", err)
			}
			if diff := cmp.Diff(tc.want, got, cmpopts.IgnoreUnexported()); diff != "" {
				t.Errorf("Post mistmatch (-want, +got): %s", diff)
			}
		})
	}
}
