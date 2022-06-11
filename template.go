package site

import (
	"html/template"
	"io"
)

type renderer struct {
	t *template.Template
}

func (r *renderer) RenderPostList(w io.Writer, posts []Post) error {
	return nil
}
