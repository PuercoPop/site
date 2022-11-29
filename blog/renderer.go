package blog

import "html/template"

type Renderer struct {
	// templates
	indextmpl   *template.Template
	postTmpl    *template.Template
	tagListTmpl *template.Template
}
