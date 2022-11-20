package finsta

import (
	"embed"
	"html/template"
	"io"

	"github.com/PuercoPop/site/blog"
)

//go:embed template/*.tmpl
var FSTemplates embed.FS

type renderer struct {
	t *template.Template
}

func (r *renderer) RenderPostList(w io.Writer, posts []blog.Post) error {
	return nil
}
