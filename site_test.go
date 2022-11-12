package site

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type FakeWWW struct{}
type FakeBlog struct{}
type FakeFinsta struct{}

func (www *FakeWWW) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("reached fakewww"))
}

func (blog *FakeBlog) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("reached fakeblog"))
}

func (finsta *FakeFinsta) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("reached fakefinsta"))
}

func TestHostRoutuing(t *testing.T) {
	tt := []struct {
		desc string
		url  string
		want string
	}{
		{
			desc: "The blog subdomain should be redirected to the blog",
			url:  "https://blog.puercopop.com",
			want: "reached fakeblog",
		},
		{
			desc: "The finsta subdomain should be redirected to finsta",
			url:  "https://finsta.puercopop.com",
			want: "reached fakefinsta",
		},
		{
			desc: "The www subdomain should be redirected to the main site",
			url:  "https://www.puercopop.com",
			want: "reached fakewww",
		},
		{
			desc: "The tld should be redirected to the main site",
			url:  "https://puercopop.com",
			want: "reached fakewww",
		},
	}
	srv := site{WWW: &FakeWWW{}, Blog: &FakeBlog{}, Finsta: &FakeFinsta{}}
	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.url, nil)
			w := httptest.NewRecorder()
			srv.ServeHTTP(w, req)
			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)
			if diff := cmp.Diff(string(body), tc.want); diff != "" {
				t.Errorf("Response mismatch (-want, +got): %s", diff)
			}
		})
	}
}
