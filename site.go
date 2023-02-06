package site

import (
	"embed"
)

//go:embed content/posts/*.md
var FSBlog embed.FS

// Add an html/template here

// func (srv *swiki) PageHandlerFunc() http.HandlerFunc {

// }
