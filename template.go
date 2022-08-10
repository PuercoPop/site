package site

import (
	"github.com/PuercoPop/site/blog"
	"html/template"
	"io"
)

type renderer struct {
	t *template.Template
}

func (r *renderer) RenderPostList(w io.Writer, posts []blog.Post) error {
	return nil
}
