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
	w.Write([]byte("reached finsta"))
}

func TestMainSite(t *testing.T) {
	// fake request from blog.puercopop.com
	// finsta.puercopop.com
	// WWW.puercopop.com
	// tt := []struct{
	// 	desc string
	// 	url string
	// 	want string
	// }
	srv := site{WWW: &FakeWWW{}, Blog: &FakeBlog{}, Finsta: &FakeFinsta{}}
	req := httptest.NewRequest(http.MethodGet, "https://blog.puercopop.com", nil)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)
	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	if diff := cmp.Diff(string(body), "reached fakeblog"); diff != "" {
		t.Errorf("Response mismatch (-want, +got): %s", diff)
	}
}
