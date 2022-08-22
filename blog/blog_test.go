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
			// TODO(javier): Move to reading from file
			Content: bytes.NewBufferString("<h1>Some title</h1>\n<h2>en, testing</h2>\n<h2>2022-30-3</h2>\n<h1>Preface</h1>\n<p>Here is some content</p>\n<h1>Another header</h1>\n"),
		},
	}}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ReadPost(tc.path)
			if err != nil {
				t.Errorf("Could not read post successfully. %s", err)
			}
			if diff := cmp.Diff(tc.want, got, cmpopts.IgnoreFields(Post{}, "Content")); diff != "" {
				t.Errorf("Post mistmatch (-want, +got): %s", diff)
			}
			if diff := cmp.Diff(tc.want.Content.String(), got.Content.String()); diff != "" {
				t.Errorf("Post Content mismatch (-want, +got): %s", diff)
			}
		})
	}
}
