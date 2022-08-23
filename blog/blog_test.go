package blog

import (
	"bytes"
	"embed"
	"io/fs"
	"testing"
	"time"

	"cloud.google.com/go/civil"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

//go:embed testdata/*.md
var FSBlog embed.FS

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
			Content: bytes.NewBufferString("<h1>Some title</h1>\n<h2>en, testing</h2>\n<h2>2022-3-30</h2>\n<h1>Preface</h1>\n<p>Here is some content</p>\n<h1>Another header</h1>\n"),
		},
	}}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ReadPost(tc.path, FSBlog)
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

func TestSite(t *testing.T) {
	blogdir, err := fs.Sub(FSBlog, "testdata")
	if err != nil {
		t.Fatalf("Could not access subdirectory: %s", err)
	}
	site := New(blogdir)
	// assert tags
	// The tags are en es testing and blog.
	var tags []string
	for t := range site.ByTag {
		tags = append(tags, t)
	}
	less := func(a, b string) bool { return a < b }
	if diff := cmp.Diff([]string{"en", "es", "testing", "blog"}, tags, cmpopts.SortSlices(less)); diff != "" {
		t.Errorf("Tag list mismatch (-want, +got): %s", diff)

	}
	// The en tag has two posts, titled 'Some title' and 'Another title'.
	var titles []string
	for _, p := range site.ByTag["en"] {
		titles = append(titles, p.Title)
	}
	if diff := cmp.Diff([]string{"Some title", "Another title"}, titles, cmpopts.SortSlices(less)); diff != "" {
		t.Errorf("en tag content mismatch (-want, +got): %s", diff)
	}
	// assert posts in date
	// 2022-30-3 TODO: switch to YYYY-M-D
	// There are three posts under date, titled 'Some title', 'Another
	// title' and 'Yet another title'.
}
