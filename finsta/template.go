package finsta

import (
	"html/template"
	"io"

	"github.com/PuercoPop/site/blog"
)

type renderer struct {
	t *template.Template
}

func (r *renderer) RenderPostList(w io.Writer, posts []blog.Post) error {
	return nil
}
