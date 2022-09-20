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
	less := func(a, b string) bool { return a < b }
	t.Run("A list of tags is assembled", func(t *testing.T) {
		var tags []string
		for t := range site.ByTag {
			tags = append(tags, t)
		}
		if diff := cmp.Diff([]string{"en", "es", "testing", "blog"}, tags, cmpopts.SortSlices(less)); diff != "" {
			t.Errorf("Tag list mismatch (-want, +got): %s", diff)

		}
	})
	t.Run("Retrieve the posts by tag", func(t *testing.T) {

		var titles []string
		for _, p := range site.ByTag["en"] {
			titles = append(titles, p.Title)
		}
		diff := cmp.Diff([]string{"Some title", "Another title"}, titles,
			cmpopts.SortSlices(less))
		if diff != "" {
			t.Errorf("en tag content mismatch (-want, +got): %s", diff)
		}
	})
	// Test ByDate index
	t.Run("Retrieve the posts by published date", func(t *testing.T) {
		var titles []string
		for _, p := range site.ByDate[civil.Date{Year: 2022, Month: time.March, Day: 30}] {
			titles = append(titles, p.Title)
		}
		diff := cmp.Diff([]string{"Some title", "Another title", "Yet another title"}, titles,
			cmpopts.SortSlices(less))
		if diff != "" {
			t.Errorf("2022-3-30 date content mismatch (-want, +got): %s", diff)
		}
	})
	// Test BySlug index
	t.Run("Retrieve a post by slug", func(t *testing.T) {
		slug := "yet-another-title"
		p := site.BySlug[slug]
		if p == nil {
			t.Fatalf("Byslug index mismatch. No post found under '%s'", slug)
		}
		if diff := cmp.Diff(p.Title, "Yet another title"); diff != "" {
			t.Errorf("BySlug index mismatch (-want, +got): %s", diff)
		}

	})
}

func TestTagList(t *testing.T) {
	tt := []struct {
		description string
		tagIndex    map[string][]*Post
		want        []tag
	}{{
		description: "With no posts",
		tagIndex:    map[string][]*Post{},
		want:        []tag{},
	},
		{
			description: "With some posts",
			tagIndex: map[string][]*Post{
				"en": []*Post{{Title: "hello"}},
				"es": []*Post{{Title: "hola"}, {Title: "mundo"}},
			},
			want: []tag{{name: "en", count: 1}, {name: "es", count: 2}},
		},
	}
	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			got := tagList(tc.tagIndex)
			if diff := cmp.Diff(tc.want, got, cmp.AllowUnexported(tag{})); diff != "" {
				t.Errorf("tag list did not match (-want, +got): %s", diff)
			}
		})
	}
}
